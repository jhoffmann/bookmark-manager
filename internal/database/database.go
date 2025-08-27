// Package database provides SQLite3 database connectivity and ORM functionality
// using GORM.
package database

import (
	"fmt"
	"log"

	"github.com/jhoffmann/bookmark-manager/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB represents the database interface
type DB interface {
	// Connection management
	Close() error
	Ping() error
	// Internal method for testing
	GetDB() *gorm.DB
}

// Database wraps the GORM database instance
type Database struct {
	db *gorm.DB
}

// Logger interface for dependency injection in tests
type Logger interface {
	Printf(format string, v ...interface{})
}

// DefaultLogger implements Logger interface
type DefaultLogger struct{}

// Printf implements the Logger interface
func (l *DefaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// NewDatabase creates a new database connection with the provided configuration
func NewDatabase(cfg *config.Config) (DB, error) {
	return NewDatabaseWithLogger(cfg, &DefaultLogger{})
}

// NewDatabaseWithLogger creates a new database connection with custom logger (for testing)
func NewDatabaseWithLogger(cfg *config.Config, appLogger Logger) (DB, error) {
	// Configure GORM logger level
	logLevel := logger.Warn
	switch cfg.GetLogLevel() {
	case "info":
		logLevel = logger.Info
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	}

	db, err := gorm.Open(sqlite.Open(cfg.GetDatabasePath()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	database := &Database{db: db}

	if err := database.migrateWithLogger(appLogger); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}

// migrate runs auto-migration for all models
func (d *Database) migrate() error {
	return d.migrateWithLogger(&DefaultLogger{})
}

// migrateWithLogger runs auto-migration for all models with custom logger
func (d *Database) migrateWithLogger(appLogger Logger) error {
	// Auto-migrate bookmark model
	if err := d.db.AutoMigrate(&BookmarkModel{}); err != nil {
		return fmt.Errorf("failed to auto-migrate bookmark table: %w", err)
	}
	appLogger.Printf("Database migration completed successfully")
	return nil
}

// BookmarkModel represents the bookmark table structure for migration
// This is a minimal model definition for auto-migration purposes
type BookmarkModel struct {
	ID          uint    `gorm:"primaryKey"`
	Folder      string  `gorm:"not null"`
	DateCreated string  `gorm:"type:datetime"`
	Category    string  `gorm:"type:varchar(50);default:'default'"`
	CreatedAt   string  `gorm:"type:datetime"`
	UpdatedAt   string  `gorm:"type:datetime"`
	DeletedAt   *string `gorm:"index;type:datetime"`
}

// TableName specifies the table name for the bookmark model
func (BookmarkModel) TableName() string {
	return "bookmarks"
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// Ping checks if the database connection is alive
func (d *Database) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// GetDB returns the underlying GORM database instance (for testing)
func (d *Database) GetDB() *gorm.DB {
	return d.db
}
