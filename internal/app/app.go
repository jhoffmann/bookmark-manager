// Package app provides centralized application initialization and setup.
package app

import (
	"fmt"
	"os"

	"github.com/jhoffmann/bookmark-manager/internal/bookmark"
	"github.com/jhoffmann/bookmark-manager/internal/config"
	"github.com/jhoffmann/bookmark-manager/internal/database"
	"github.com/jhoffmann/bookmark-manager/internal/tui/styles"
)

// App holds the database connection and bookmark service
type App struct {
	DB      database.DB
	Service *bookmark.Service
	Config  *config.Config
}

// Close closes the database connection
func (a *App) Close() error {
	return a.DB.Close()
}

// Initialize loads configuration, initializes database, and returns an App instance
// This centralizes all the repetitive setup code from the command files
func Initialize() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize bookmark service
	service := bookmark.NewService(db)

	return &App{
		DB:      db,
		Service: service,
		Config:  cfg,
	}, nil
}

// InitializeOrExit initializes the app and exits on error with styled messages
// This is a convenience function for commands that should exit on initialization failure
func InitializeOrExit() *App {
	app, err := Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n",
			styles.ErrorMessage.Render("✗"), err)
		os.Exit(1)
	}
	return app
}

// MustInitialize is like Initialize but panics on error (useful for tests)
func MustInitialize() *App {
	app, err := Initialize()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize app: %v", err))
	}
	return app
}

// InitializeWithConfig initializes with a specific configuration (useful for testing)
func InitializeWithConfig(cfg *config.Config) (*App, error) {
	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize bookmark service
	service := bookmark.NewService(db)

	return &App{
		DB:      db,
		Service: service,
		Config:  cfg,
	}, nil
}

// InitializeInMemory creates an in-memory database for testing
func InitializeInMemory() (*App, error) {
	cfg := &config.Config{
		DatabasePath: ":memory:",
		LogLevel:     "silent",
	}

	return InitializeWithConfig(cfg)
}
