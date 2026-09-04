package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("fail to get home dir: %w", err)
	}

	fullPath := dir + "/" + configFileName

	return fullPath, nil
}

func write(c Config) error {

	filePath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	marshedData, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("Failed to convert data to json: %w", err)
	}

	if err := os.WriteFile(filePath, marshedData, 0644); err != nil {
		return fmt.Errorf("Fialed to update config file: %w", err)
	}

	return nil

}

func GetConfig() (*Config, error) {

	filePath, err := getConfigFilePath()

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("fail to read file: %w", err)
	}

	var newConfig Config

	if err := json.Unmarshal(data, &newConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &newConfig, nil

}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	if err := write(*c); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetUser() string {

	username := c.CurrentUserName

	return username
}
