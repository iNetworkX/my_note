package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configFileName = ".note-cli-config.json"
)

type Config struct {
	NotesDirectory string `json:"notes_directory"`
}

type ConfigManager struct {
	configPath string
}

func NewConfigManager() *ConfigManager {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, configFileName)
	
	return &ConfigManager{
		configPath: configPath,
	}
}

func (c *ConfigManager) Load() (*Config, error) {
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return nil, err
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

func (c *ConfigManager) Save(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(c.configPath, data, 0600)
}

func (c *ConfigManager) Exists() bool {
	_, err := os.Stat(c.configPath)
	return err == nil
}

func (c *ConfigManager) GetNotesDirectory() (string, error) {
	config, err := c.Load()
	if err != nil {
		return "", err
	}
	
	return config.NotesDirectory, nil
}

func (c *ConfigManager) SetupNotesDirectory() (string, error) {
	fmt.Println("\nFirst time setup - Please choose where to store your encrypted notes.")
	fmt.Println("Examples: ~/Dropbox/, ~/Documents/notes/, ~/.my_notes/")
	fmt.Print("Enter directory path: ")
	
	var dirPath string
	if _, err := fmt.Scanln(&dirPath); err != nil {
		return "", fmt.Errorf("failed to read directory path: %v", err)
	}
	
	// Expand tilde to home directory
	if len(dirPath) >= 2 && dirPath[:2] == "~/" {
		homeDir, _ := os.UserHomeDir()
		dirPath = filepath.Join(homeDir, dirPath[2:])
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	
	// Save config
	config := &Config{
		NotesDirectory: dirPath,
	}
	
	if err := c.Save(config); err != nil {
		return "", fmt.Errorf("failed to save config: %v", err)
	}
	
	fmt.Printf("Notes will be stored in: %s\n", dirPath)
	return dirPath, nil
}