package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	JSON_NAME = ".gatorconfig.json"
)

type Config struct {
	DB_URL            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

// Returns ~/JSON_NAME filepath
func getConfigFilePath() (string, error) {
	// Get HOME dir
	homeDir, errDir := os.UserHomeDir()
	if errDir != nil {
		return "", fmt.Errorf("error when trying to get HOME dir: %w", errDir)
	}

	return homeDir + "/" + JSON_NAME, nil
}

func write(cfg Config) error {
	configPath, errConfigPath := getConfigFilePath()
	if errConfigPath != nil {
		return fmt.Errorf("error when trying to get config file path: %w", errConfigPath)
	}
	configData, errMarshal := json.Marshal(cfg)
	if errMarshal != nil {
		return fmt.Errorf("error when trying to marshal Config struct: %w", errMarshal)
	}
	errWrite := os.WriteFile(configPath, configData, 0777)
	if errWrite != nil {
		return fmt.Errorf("error when trying to write to config file: %w", errWrite)
	}

	return nil
}

// Reads config file and returns it into Config struct
func Read() (Config, error) {
	// Get HOME dir
	configPath, errConfigPath := getConfigFilePath()
	if errConfigPath!= nil {
		return Config{}, fmt.Errorf("error when trying to get config file path: %w", errConfigPath)
	}
	// Read ~/JSON_NAME file
	fileData, errRead := os.ReadFile(configPath)
	if errRead != nil {
		return Config{}, fmt.Errorf("error when trying to read JSON file: %w", errRead)
	}
	// Unmarshal read file
	var readConfig Config
	errUnmarshal := json.Unmarshal(fileData, &readConfig)
	if errUnmarshal != nil {
		return Config{}, fmt.Errorf("error when trying to unmarshal JSON data: %w", errUnmarshal)
	}

	return readConfig, nil
}

func (c *Config) SetUser(currentUserName string) error {
	// Set config field
	c.Current_user_name = currentUserName
	// Write to config file
	errWrite := write(*c)
	if errWrite != nil {
		return fmt.Errorf("error when trying to set user: %w", errWrite)
	}

	return nil
}
