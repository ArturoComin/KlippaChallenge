package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

// Config represents the configuration settings for your CLI tool.
type Config struct {
	APIKey            string `json:"api_key"`
	Template          string `json:"template"`
	PDFTextExtraction string `json:"pdf_text_extraction"`
	// Add other configuration options here
}

var globalConfig *Config

// LoadConfig loads the global configuration from a config file (e.g., ~/.my-cli-tool/config.json).
func LoadConfig() {
	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		return
	}
	homeDir := currentUser.HomeDir

	// Define the path to the config file
	configFilePath := filepath.Join(homeDir, ".my-cli-tool", "config.json")

	// Check if the config file exists
	if _, err := os.Stat(configFilePath); err != nil {
		// Config file doesn't exist, use default values or initialize as needed
		globalConfig = &Config{
			APIKey:            "",
			Template:          "",
			PDFTextExtraction: "fast",
			// Initialize other options as needed
		}
		return
	}

	// Read and parse the config file
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&globalConfig); err != nil {
		return
	}
}

// GetGlobalConfig returns the loaded global configuration.
func GetGlobalConfig() *Config {
	return globalConfig
}
