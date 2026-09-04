package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JiveshL-KDK/go-gator/internal/database"
	"github.com/JiveshL-KDK/go-gator/internal/state"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

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
