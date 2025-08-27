// Package config provides configuration management for the bookmark manager application.
// It handles loading configuration from environment variables and provides cross-platform
// default paths for application data.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds application configuration options
type Config struct {
	DatabasePath string `envconfig:"BM_DATABASE" json:"database_path"`
	LogLevel     string `envconfig:"BM_LOGLEVEL" json:"log_level"`
}

// Load loads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	config := &Config{
		LogLevel: "warn", // Default log level
	}

	// Load from environment variables
	if dbPath := os.Getenv("BM_DATABASE"); dbPath != "" {
		config.DatabasePath = dbPath
	} else {
		// Use default path in user's config directory
		defaultPath, err := getDefaultDatabasePath()
		if err != nil {
			return nil, fmt.Errorf("failed to get default database path: %w", err)
		}
		config.DatabasePath = defaultPath
	}

	if logLevel := os.Getenv("BM_LOGLEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	return config, nil
}

// getDefaultDatabasePath returns the default database path for the current platform
// Linux/Unix: ~/.config/bookmark-manager/bookmarks.db
// macOS: ~/Library/Application Support/bookmark-manager/bookmarks.db
// Windows: %APPDATA%/bookmark-manager/bookmarks.db
func getDefaultDatabasePath() (string, error) {
	var configDir string
	var err error

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support")
	default: // Linux and other Unix-like systems
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			configDir = filepath.Join(homeDir, ".config")
		}
	}

	appDir := filepath.Join(configDir, "bookmark-manager")
	dbPath := filepath.Join(appDir, "bookmarks.db")

	// Ensure the directory exists
	if err = os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", appDir, err)
	}

	return dbPath, nil
}

// GetDatabasePath returns the configured database path
func (c *Config) GetDatabasePath() string {
	return c.DatabasePath
}

// GetLogLevel returns the configured log level
func (c *Config) GetLogLevel() string {
	return c.LogLevel
}
