package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Set default name and content for .JSON config file
const (
	JSON_NAME    = ".rssconfig.json"
	INITIAL_JSON = `{"db_url": "postgres://postgres:postgres@localhost:5432/rss?sslmode=disable"}`
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Returns ~/JSON_NAME filepath
func getConfigFilePath() (string, error) {

	// Get HOME dir
	homeDir, errDir := os.UserHomeDir()
	if errDir != nil {
		return "", fmt.Errorf("couldn't get HOME dir: %w", errDir)
	}

	// Construct full path (OS-agnostic)
	fullPath := filepath.Join(homeDir, JSON_NAME)

	// Return normally
	return fullPath, nil
}

// Reads config file and returns it into a Config struct
func Read() (Config, error) {

	// Get HOME dir
	configPath, errConfigPath := getConfigFilePath()
	if errConfigPath != nil {
		return Config{}, fmt.Errorf("couldn't get config file path: %w", errConfigPath)
	}

	// Read ~/JSON_NAME file, writes it if not found
	_, errRead := os.ReadFile(configPath)
	if errRead != nil && errRead.Error() == fmt.Sprintf("open %s: no such file or directory", configPath) {
		fmt.Println(".rssconfig.json not found, writing it to HOME dir")
		errWrite := os.WriteFile(configPath, []byte(INITIAL_JSON), 0666)
		if errWrite != nil {
			return Config{}, fmt.Errorf("couldn't write initial JSON: %w", errWrite)
		}
	}

	// Read ~/JSON_NAME file again and return if an error is found
	fileData, errReadTwo := os.ReadFile(configPath)
	if errReadTwo != nil {
		return Config{}, fmt.Errorf("couldn't read initial JSON file: %w", errReadTwo)
	}

	// Unmarshal read file
	var readConfig Config
	errUnmarshal := json.Unmarshal(fileData, &readConfig)
	if errUnmarshal != nil {
		return Config{}, fmt.Errorf("couldn't unmarshal JSON data: %w", errUnmarshal)
	}

	// Return normally
	return readConfig, nil
}

// Sets username and writes Config struct to disk
func (c *Config) SetUser(currentUserName string) error {

	// Set config field
	c.CurrentUserName = currentUserName

	// Write to config file
	errWrite := write(*c)
	if errWrite != nil {
		return fmt.Errorf("couldn't set user: %w", errWrite)
	}

	// Return normally
	return nil
}

// Writes the Config struct to ~/JSON_NAME filepath
func write(cfg Config) error {

	// Get config file path
	configPath, errConfigPath := getConfigFilePath()
	if errConfigPath != nil {
		return fmt.Errorf("couldn't get config file path: %w", errConfigPath)
	}

	// Marshal config struct
	configData, errMarshal := json.Marshal(cfg)
	if errMarshal != nil {
		return fmt.Errorf("couldn't marshal Config struct: %w", errMarshal)
	}

	// Write to config file
	errWrite := os.WriteFile(configPath, configData, 0666)
	if errWrite != nil {
		return fmt.Errorf("couldn't write to config file: %w", errWrite)
	}

	// Return normally
	return nil
}