package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jhoffmann/bookmark-manager/internal/config"
)

func TestNewDatabase(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	// Test successful database creation
	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase() error = %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Error("NewDatabase() returned nil database")
	}

	// Test that the database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestNewDatabase_InvalidPath(t *testing.T) {
	cfg := &config.Config{
		DatabasePath: "/invalid/path/that/does/not/exist/test.db",
		LogLevel:     "silent",
	}

	db, err := NewDatabase(cfg)
	if err == nil {
		if db != nil {
			db.Close()
		}
		t.Error("NewDatabase() expected error for invalid path, got nil")
	}

	if db != nil {
		t.Error("NewDatabase() should return nil database on error")
	}
}

func TestDatabase_LogLevels(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		logLevel string
	}{
		{"info level", "info"},
		{"error level", "error"},
		{"silent level", "silent"},
		{"warn level", "warn"},
		{"default level", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbPath := filepath.Join(tempDir, "test_"+tt.name+".db")
			cfg := &config.Config{
				DatabasePath: dbPath,
				LogLevel:     tt.logLevel,
			}

			db, err := NewDatabase(cfg)
			if err != nil {
				t.Fatalf("NewDatabase() error = %v", err)
			}
			defer db.Close()

			if db == nil {
				t.Error("NewDatabase() returned nil database")
			}
		})
	}
}

func TestDatabase_Close(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_close.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase() error = %v", err)
	}

	// Test successful close
	err = db.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Test that we can close multiple times (should not panic)
	err = db.Close()
	if err == nil {
		t.Log("Close() on already closed database returned no error (this is acceptable)")
	}
}

func TestDatabase_Ping(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_ping.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase() error = %v", err)
	}
	defer db.Close()

	// Test successful ping
	err = db.Ping()
	if err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}

func TestDatabase_GetDB(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_getdb.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase() error = %v", err)
	}
	defer db.Close()

	gormDB := db.GetDB()
	if gormDB == nil {
		t.Error("GetDB() returned nil GORM database")
	}
}

func TestDatabase_Migrate(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_migrate.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	// Create database instance
	database, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase() error = %v", err)
	}
	defer database.Close()

	// Cast to concrete type to test migrate method directly
	db, ok := database.(*Database)
	if !ok {
		t.Fatal("Expected *Database type")
	}

	// Test migrate method directly
	err = db.migrate()
	if err != nil {
		t.Errorf("migrate() error = %v", err)
	}
}

// Integration test with actual SQLite operations
func TestDatabase_Integration(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "integration_test.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase() error = %v", err)
	}
	defer db.Close()

	// Test ping
	if err := db.Ping(); err != nil {
		t.Errorf("Ping() error = %v", err)
	}

	// Test GetDB returns valid GORM instance
	gormDB := db.GetDB()
	if gormDB == nil {
		t.Fatal("GetDB() returned nil")
	}

	// Test that we can perform a simple query
	var result int64
	if err := gormDB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		t.Errorf("Simple query failed: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected result 1, got %d", result)
	}

	// Test close
	if err := db.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}
