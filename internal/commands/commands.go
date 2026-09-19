package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/JiveshL-KDK/go-gator/internal/constants"
	"github.com/JiveshL-KDK/go-gator/internal/database"
	"github.com/JiveshL-KDK/go-gator/internal/rss"
	"github.com/JiveshL-KDK/go-gator/internal/state"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

func middlewareLoggedIn(handler func(s *state.State, cmd CommandInput) error) func(*state.State, CommandInput) error {

	return func(s *state.State, cmd CommandInput) error {
		user := s.GetUser()

		if user == "" {
			fmt.Println("Kindly login first")
			return nil
		}

		return handler(s, cmd)

	}
}

func handleLogin(s *state.State, cmd CommandInput) error {
	args := cmd.args

	if len(args) == 0 {
		fmt.Printf("Currently Logged in as: %s\n", s.GetUser())
		return nil
	}

	user := args[0]

	err := s.SetUser(user)

	if err != nil {
		return fmt.Errorf("failed to login!: %w", err)
	}

	fmt.Printf("Logged in as %s\n", user)

	return nil

}

func handleRegister(s *state.State, cmd CommandInput) error {
	args := cmd.args

	if len(args) == 0 {
		return fmt.Errorf("Kindly pass the name of the user you want to register!")
	}

	user := args[0]

	uuid, err := uuid.NewRandom()

	if err != nil {
		return fmt.Errorf("Failed to create a UUID: %w", err)
	}

	newUser := database.CreateUserParams{ID: uuid, Name: user, CreatedAt: sql.NullTime{Time: time.Now(), Valid: true}, UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true}}

	createdUser, err := s.Db.CreateUser(context.Background(), newUser)

	if err != nil {

		pqError, ok := errors.AsType[*pq.Error](err)

		if !ok {
			return fmt.Errorf("Failed to create a new user: %w", err)
		}

		if pqError.Code == pqerror.UniqueViolation {
			if pqError.Constraint == "users_name_key" {
				return fmt.Errorf("User with %v already exists!", user)
			}

			if pqError.Constraint == "users_pkey" {
				return fmt.Errorf("Watch out for any meteors in the sky, anyways try that again!")
			}

		}
		return fmt.Errorf("Failed to create the user %w", pqError)

	}

	fmt.Printf("New User: %s with Id: %s Created!\n", createdUser.Name, createdUser.ID.String())

	return nil
}

func listUsers(s *state.State, cmd CommandInput) error {

	users, err := s.Db.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("Failed to get users: %w", err)
	}

	for _, user := range users {
		fmt.Println(user.Name)
	}

	return nil

}

func resetAllUsers(s *state.State, cmd CommandInput) error {

	err := s.Db.RemoveAllUsers(context.Background())

	if err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}

	return nil
}

func agg(s *state.State, cmd CommandInput) error {

	timeBetweenRequests, err := time.ParseDuration(constants.TIME_BW_REQ)

	if err != nil {
		return fmt.Errorf("Failed to parse the duration: %w", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		rss.ScrapeFeeds(s, context.Background())
	}
}

func addFeed(s *state.State, cmd CommandInput) error {

	args := cmd.args

	if len(args) == 0 {
		return fmt.Errorf("Kindly provide the name and url of the feed you want to add")
	}

	if len(args) == 1 {
		return fmt.Errorf("Kinldy provide both name and url")
	}

	name := args[0]
	urlInput := args[1]

	if len(name) > constants.MAX_FEED_TITLE_LIMIT {
		return fmt.Errorf("Kindly provide the len of name less than %d chars", constants.MAX_FEED_TITLE_LIMIT)
	}

	if _, err := url.Parse(urlInput); err != nil {
		return fmt.Errorf("Kindly provide a valid feed url")
	}

	currentUser := s.GetUser()
	currentUserRecord, err := s.Db.GetUser(context.Background(), currentUser)

	if err != nil {
		return fmt.Errorf("Failed to get user record: %w", err)
	}

	if ok, err := s.Db.UserHasFeed(context.Background(), database.UserHasFeedParams{UserID: currentUserRecord.ID, Url: urlInput}); err != nil {
		return fmt.Errorf("Fail to check user feeds, %w", err)
	} else if ok {
		return fmt.Errorf("Feed already registerd!")
	}

	id, err := uuid.NewUUID()

	if err != nil {
		return fmt.Errorf("Failed to make new UUID: %w", err)
	}

	_, err = s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        id,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      name,
		Url:       urlInput,
		UserID:    currentUserRecord.ID,
	})

	if err != nil {
		return fmt.Errorf("faied to create user feed")
	}

	feedFollowKey, err := uuid.NewRandom()

	if err != nil {
		return fmt.Errorf("Failed to create a primary key: %w", err)
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        feedFollowKey,
		UserID:    currentUserRecord.ID,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		FeedID:    id,
	})

	if err != nil {
		return fmt.Errorf("Failed to create a feed: %w", err)
	}

	return nil

}

func listFeeds(s *state.State, cmd CommandInput) error {
	currentUserName := s.GetUser()

	user, err := s.Db.GetUser(context.Background(), currentUserName)

	if err != nil {
		return fmt.Errorf("Failed to get user properly:- %w", err)
	}

	feeds, err := s.Db.GetFeedsOfUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("Failed to get feeds of this user: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds for you currently!")
		return nil
	}

	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println()
	}

	return nil

}

func follow(s *state.State, cmd CommandInput) error {

	args := cmd.args

	if len(args) == 0 {
		return fmt.Errorf("Kinldy provide the url to add")
	}

	urlInput := args[0]

	if _, err := url.Parse(urlInput); err != nil {
		return fmt.Errorf("kindly provide a valid url: %w", err)
	}

	feed, err := s.Db.GetFeedByURL(context.Background(), urlInput)

	if err != nil {
		return fmt.Errorf("Failed to get feed: %w", err)
	}

	user, err := s.Db.GetUser(context.Background(), s.GetUser())

	if err != nil {
		return fmt.Errorf("Failed to get current logged in user: %w", err)
	}

	feedFollowKey, err := uuid.NewRandom()

	if err != nil {
		return fmt.Errorf("Failed to create a primary key: %w", err)
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        feedFollowKey,
		UserID:    user.ID,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Failed to create a feed: %w", err)
	}

	fmt.Println(feedFollow)

	return nil

}

func following(s *state.State, cmd CommandInput) error {
	userName := s.GetUser()

	user, err := s.Db.GetUser(context.Background(), userName)

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	follows, err := s.Db.GetFeedFollowsOfUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("failed to get follows of the user: %w", err)
	}

	for _, follow := range follows {
		fmt.Println(follow.FeedName)
		fmt.Println(follow.UserName)
		fmt.Println()
	}

	return nil
}

func unfollow(s *state.State, cmd CommandInput) error {

	userName := s.GetUser()

	user, err := s.Db.GetUser(context.Background(), userName)

	if err != nil {
		return fmt.Errorf("Fialed to get user record: %w", err)
	}

	args := cmd.args

	if len(args) == 0 {
		return fmt.Errorf("Kindly provide the feed name")
	}

	feed, err := s.Db.GetFeedByURL(context.Background(), args[0])

	if err != nil {
		return fmt.Errorf("Failed to get feed: %w", err)
	}

	if err := s.Db.RemoveFeedFollow(context.Background(), database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("failed to delete the record: %w", err)
	}

	return nil

}

func browse(s *state.State, cmd CommandInput) error {
	currentUser := s.GetUser()

	currentUserRecord, err := s.Db.GetUser(context.Background(), currentUser)

	if err != nil {
		return fmt.Errorf("Failed to get current user: %w", err)
	}

	limit := 2

	if len(cmd.args) > 0 {
		conv, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("Kinldy provide a valid argument: %w", err)
		}

		limit = conv
	}

	currentPosts, err := s.Db.GetPostsForAUser(context.Background(), database.GetPostsForAUserParams{
		UserID: currentUserRecord.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return fmt.Errorf("Failed to get posts: %w", err)
	}

	for _, post := range currentPosts {
		fmt.Println(post.Title)
		fmt.Println(post.Description.String)
		fmt.Println()
	}

	return nil
}
