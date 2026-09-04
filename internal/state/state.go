package state

import (
	"database/sql"
	"fmt"

	"github.com/JiveshL-KDK/go-gator/internal/config"
	"github.com/JiveshL-KDK/go-gator/internal/database"

	_ "github.com/lib/pq"
)

type State struct {
	config *config.Config
	Db     *database.Queries
}

func (s *State) initalizeConfig() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("Filed to initalize config for state: %w", err)
	}

	s.config = cfg

	return nil
}

func (s *State) initializeDB() error {
	dbUrl := s.config.DbUrl

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		return fmt.Errorf("Fail to open db conneciton: %w", err)
	}

	dbQueries := database.New(db)
	s.Db = dbQueries

	return nil
}

func (s *State) SetUser(username string) error {
	if err := s.config.SetUser(username); err != nil {
		return err
	}

	return nil
}

func (s *State) GetUser() string {
	username := s.config.GetUser()
	return username
}

func InitialzeState() (*State, error) {

	newState := State{}
	if err := newState.initalizeConfig(); err != nil {
		return nil, fmt.Errorf("Failed to initalize state: %w", err)
	}

	if err := newState.initializeDB(); err != nil {
		return nil, fmt.Errorf("Failed to initialize state: %w", err)
	}

	return &newState, nil

}
