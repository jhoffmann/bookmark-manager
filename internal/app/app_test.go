package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jhoffmann/bookmark-manager/internal/bookmark"
	"github.com/jhoffmann/bookmark-manager/internal/config"
)

func TestInitialize(t *testing.T) {
	// Set temporary database path to avoid conflicts
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_app.db")

	originalDB := os.Getenv("BM_DATABASE")
	os.Setenv("BM_DATABASE", dbPath)
	defer os.Setenv("BM_DATABASE", originalDB)

	app, err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	defer app.Close()

	if app.DB == nil {
		t.Error("Initialize() app.DB is nil")
	}

	if app.Service == nil {
		t.Error("Initialize() app.Service is nil")
	}

	if app.Config == nil {
		t.Error("Initialize() app.Config is nil")
	}

	// Verify the config path is correct
	if app.Config.GetDatabasePath() != dbPath {
		t.Errorf("Expected database path %s, got %s", dbPath, app.Config.GetDatabasePath())
	}
}

func TestInitializeWithConfig(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_with_config.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	app, err := InitializeWithConfig(cfg)
	if err != nil {
		t.Fatalf("InitializeWithConfig() error = %v", err)
	}
	defer app.Close()

	if app.DB == nil {
		t.Error("InitializeWithConfig() app.DB is nil")
	}

	if app.Service == nil {
		t.Error("InitializeWithConfig() app.Service is nil")
	}

	if app.Config != cfg {
		t.Error("InitializeWithConfig() app.Config is not the same as provided config")
	}

	// Test that we can use the service
	testBookmark := &bookmark.Bookmark{
		Folder:   "/test/folder",
		Category: "test",
	}

	if err := testBookmark.Save(app.Service); err != nil {
		t.Fatalf("Failed to save test bookmark: %v", err)
	}

	if testBookmark.ID == 0 {
		t.Error("Expected bookmark ID to be set after save")
	}
}

func TestInitializeInMemory(t *testing.T) {
	app, err := InitializeInMemory()
	if err != nil {
		t.Fatalf("InitializeInMemory() error = %v", err)
	}
	defer app.Close()

	if app.DB == nil {
		t.Error("InitializeInMemory() app.DB is nil")
	}

	if app.Service == nil {
		t.Error("InitializeInMemory() app.Service is nil")
	}

	if app.Config == nil {
		t.Error("InitializeInMemory() app.Config is nil")
	}

	// Verify it's using in-memory database
	if app.Config.GetDatabasePath() != ":memory:" {
		t.Errorf("Expected in-memory database path ':memory:', got %s", app.Config.GetDatabasePath())
	}

	// Test that we can use the service with in-memory database
	testBookmark := &bookmark.Bookmark{
		Folder:   "/test/memory",
		Category: "memory-test",
	}

	if err := testBookmark.Save(app.Service); err != nil {
		t.Fatalf("Failed to save test bookmark: %v", err)
	}

	// Verify the bookmark was saved
	retrieved, err := bookmark.GetByID(app.Service, testBookmark.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve test bookmark: %v", err)
	}

	if retrieved.Folder != testBookmark.Folder {
		t.Errorf("Expected folder %s, got %s", testBookmark.Folder, retrieved.Folder)
	}
}

func TestMustInitialize(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_must_init.db")

	originalDB := os.Getenv("BM_DATABASE")
	os.Setenv("BM_DATABASE", dbPath)
	defer os.Setenv("BM_DATABASE", originalDB)

	// This should not panic with valid config
	app := MustInitialize()
	defer app.Close()

	if app == nil {
		t.Error("MustInitialize() returned nil")
	}
}

func TestMustInitializePanic(t *testing.T) {
	// Set an invalid database path to force an error
	originalDB := os.Getenv("BM_DATABASE")
	os.Setenv("BM_DATABASE", "/invalid/path/that/does/not/exist/test.db")
	defer os.Setenv("BM_DATABASE", originalDB)

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustInitialize() should have panicked with invalid config")
		}
	}()

	MustInitialize()
}

func TestAppClose(t *testing.T) {
	app, err := InitializeInMemory()
	if err != nil {
		t.Fatalf("InitializeInMemory() error = %v", err)
	}

	// Test that Close() doesn't return an error
	if err := app.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Test that we can call Close() multiple times
	if err := app.Close(); err == nil {
		t.Log("Close() on already closed app returned no error (acceptable)")
	}
}
