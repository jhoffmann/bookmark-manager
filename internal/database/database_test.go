package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhoffmann/bookmark-manager/internal/config"
)

// TestLogger implements Logger interface for testing
type TestLogger struct {
	messages []string
}

func (l *TestLogger) Printf(format string, v ...interface{}) {
	message := format
	if len(v) > 0 {
		// Simple format replacement for testing
		for range v {
			message = strings.Replace(message, "%v", "%s", 1)
		}
	}
	l.messages = append(l.messages, message)
}

func (l *TestLogger) GetMessages() []string {
	return l.messages
}

func (l *TestLogger) Reset() {
	l.messages = nil
}

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

func TestNewDatabaseWithLogger(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_with_logger.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	testLogger := &TestLogger{}

	db, err := NewDatabaseWithLogger(cfg, testLogger)
	if err != nil {
		t.Fatalf("NewDatabaseWithLogger() error = %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Error("NewDatabaseWithLogger() returned nil database")
	}

	// Check that migration message was logged
	messages := testLogger.GetMessages()
	if len(messages) == 0 {
		t.Error("Expected migration log message, got none")
	}

	found := false
	for _, msg := range messages {
		if strings.Contains(msg, "Database migration completed successfully") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected migration success message, got messages: %v", messages)
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

	testLogger := &TestLogger{}

	// Create database instance
	database, err := NewDatabaseWithLogger(cfg, testLogger)
	if err != nil {
		t.Fatalf("NewDatabaseWithLogger() error = %v", err)
	}
	defer database.Close()

	// Cast to concrete type to test migrate method directly
	db, ok := database.(*Database)
	if !ok {
		t.Fatal("Expected *Database type")
	}

	// Test migrate method directly
	testLogger.Reset()
	err = db.migrate()
	if err != nil {
		t.Errorf("migrate() error = %v", err)
	}

	// The migrate() method uses DefaultLogger, not our testLogger
	// So we can't capture its messages. Just verify it doesn't error.
	// The message logging is tested in the migrateWithLogger test.
}

func TestDatabase_MigrateWithLogger(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_migrate_logger.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	database, err := NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase() error = %v", err)
	}
	defer database.Close()

	// Cast to concrete type to test migrateWithLogger method directly
	db, ok := database.(*Database)
	if !ok {
		t.Fatal("Expected *Database type")
	}

	testLogger := &TestLogger{}
	err = db.migrateWithLogger(testLogger)
	if err != nil {
		t.Errorf("migrateWithLogger() error = %v", err)
	}

	// Check that migration message was logged
	messages := testLogger.GetMessages()
	if len(messages) == 0 {
		t.Error("Expected migration log message, got none")
	}

	found := false
	for _, msg := range messages {
		if strings.Contains(msg, "Database migration completed successfully") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected migration success message, got messages: %v", messages)
	}
}

func TestDefaultLogger(t *testing.T) {
	logger := &DefaultLogger{}

	// This test just ensures the logger doesn't panic
	// We can't easily test the actual log output without capturing stdout
	logger.Printf("Test message: %s", "test")
	logger.Printf("Test message without args")
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
