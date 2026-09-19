package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/JiveshL-KDK/go-gator/internal/database"
	"github.com/JiveshL-KDK/go-gator/internal/state"
	"github.com/google/uuid"
)

func decodeStrings(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)

	if err != nil {
		return nil, fmt.Errorf("Can't make req: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Can't do req: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unsucessful req: %w", err)
	}

	var feed RSSFeed
	if err := xml.NewDecoder(res.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("Failed to parse res: %w", err)
	}

	decodeStrings(&feed)

	return &feed, nil

}

func toCreatePostsParams(posts []database.CreateNewPostParams) database.CreatePostsParams {
	p := database.CreatePostsParams{}

	for _, post := range posts {
		p.Ids = append(p.Ids, post.ID)
		p.CreatedAts = append(p.CreatedAts, post.CreatedAt.Time)
		p.UpdatedAts = append(p.UpdatedAts, post.UpdatedAt.Time)
		p.PublishedAts = append(p.PublishedAts, post.PublishedAt.Time)
		p.Titles = append(p.Titles, post.Title)
		p.Urls = append(p.Urls, post.Url)
		p.Descriptions = append(p.Descriptions, post.Description.String)
		p.FeedIds = append(p.FeedIds, post.FeedID)
	}

	return p
}

func ScrapeFeeds(s *state.State, ctx context.Context) error {

	nxtFeed, err := s.Db.GetNextFeedToFetch(ctx)

	if err != nil {
		return fmt.Errorf("Failed to get next feed %w", err)
	}

	fmt.Println("Fetching feed:- ", nxtFeed.Url)
	feed, err := FetchFeed(ctx, nxtFeed.Url)

	if err != nil {
		return fmt.Errorf("Failed to get feed from url: %w", err)
	}

	if err := s.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{ID: nxtFeed.ID}); err != nil {
		return fmt.Errorf("Failed to update feed: %w", err)
	}

	posts := []database.CreateNewPostParams{}

	for _, item := range feed.Channel.Item {

		id, err := uuid.NewRandom()

		if err != nil {
			return fmt.Errorf("fail to generate UUID for a post: %w", err)
		}

		post := database.CreateNewPostParams{
			ID:          id,
			CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			PublishedAt: sql.NullTime{Time: time.Now(), Valid: true},
			FeedID:      nxtFeed.ID,
			Title:       item.Title,
			Description: sql.NullString{String: item.Description, Valid: true},
			Url:         item.Link,
		}

		posts = append(posts, post)

	}

	if _, err := s.Db.CreatePosts(ctx, toCreatePostsParams(posts)); err != nil {
		return fmt.Errorf("failed to save posts: %w", err)
	}

	return nil

}
