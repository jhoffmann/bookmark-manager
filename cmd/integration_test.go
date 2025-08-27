package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jhoffmann/bookmark-manager/internal/app"
	"github.com/jhoffmann/bookmark-manager/internal/bookmark"
	"github.com/jhoffmann/bookmark-manager/internal/config"
)

// TestIntegration tests the end-to-end CLI functionality
func TestIntegration(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_integration.db")

	// Set environment variable for test database
	originalDB := os.Getenv("BM_DATABASE")
	os.Setenv("BM_DATABASE", dbPath)
	defer os.Setenv("BM_DATABASE", originalDB)

	// Initialize app with custom database path
	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	appInstance, err := app.InitializeWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}
	defer appInstance.Close()

	// Test adding bookmarks
	testBookmarks := []*bookmark.Bookmark{
		{Folder: "/test/work/project1", Category: "work"},
		{Folder: "/test/personal/docs", Category: "personal"},
		{Folder: "/test/temp", Category: "test"},
	}

	for _, b := range testBookmarks {
		if err := b.Save(appInstance.Service); err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}
	}

	// Test listing all bookmarks
	allBookmarks, err := bookmark.List(appInstance.Service, 0, 0)
	if err != nil {
		t.Fatalf("Failed to list bookmarks: %v", err)
	}

	if len(allBookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(allBookmarks))
	}

	// Test searching by category
	workBookmarks, err := bookmark.SearchByCategory(appInstance.Service, "work")
	if err != nil {
		t.Fatalf("Failed to search by category: %v", err)
	}

	if len(workBookmarks) != 1 {
		t.Errorf("Expected 1 work bookmark, got %d", len(workBookmarks))
	}

	if workBookmarks[0].Folder != "/test/work/project1" {
		t.Errorf("Expected work bookmark folder '/test/work/project1', got %s", workBookmarks[0].Folder)
	}

	// Test searching by folder
	tempBookmarks, err := bookmark.SearchByFolder(appInstance.Service, "temp")
	if err != nil {
		t.Fatalf("Failed to search by folder: %v", err)
	}

	if len(tempBookmarks) != 1 {
		t.Errorf("Expected 1 temp bookmark, got %d", len(tempBookmarks))
	}

	// Test deleting a bookmark
	if err := testBookmarks[0].Delete(appInstance.Service); err != nil {
		t.Fatalf("Failed to delete bookmark: %v", err)
	}

	// Verify deletion
	remainingBookmarks, err := bookmark.List(appInstance.Service, 0, 0)
	if err != nil {
		t.Fatalf("Failed to list bookmarks after deletion: %v", err)
	}

	if len(remainingBookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks after deletion, got %d", len(remainingBookmarks))
	}
}

// TestConfigurationLoading tests that configuration loads properly
func TestConfigurationLoading(t *testing.T) {
	// Test with environment variables
	originalDB := os.Getenv("BM_DATABASE")
	originalLog := os.Getenv("BM_LOGLEVEL")

	os.Setenv("BM_DATABASE", "/custom/path.db")
	os.Setenv("BM_LOGLEVEL", "info")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.GetDatabasePath() != "/custom/path.db" {
		t.Errorf("Expected database path '/custom/path.db', got %s", cfg.GetDatabasePath())
	}

	if cfg.GetLogLevel() != "info" {
		t.Errorf("Expected log level 'info', got %s", cfg.GetLogLevel())
	}

	// Restore environment variables
	os.Setenv("BM_DATABASE", originalDB)
	os.Setenv("BM_LOGLEVEL", originalLog)
}

// TestCustomCategories tests that custom user-defined categories work properly
func TestCustomCategories(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_categories.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	// Initialize app with config
	appInstance, err := app.InitializeWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize app: %v", err)
	}
	defer appInstance.Close()

	// Test bookmarks with custom categories
	customCategories := []string{
		"urgent-projects",
		"learning-materials",
		"entertainment",
		"side-projects",
		"archive",
	}

	for i, category := range customCategories {
		bookmark := &bookmark.Bookmark{
			Folder:   filepath.Join("/test", category, "folder"+string(rune('0'+i))),
			Category: bookmark.CategoryType(category),
		}

		if err := bookmark.Save(appInstance.Service); err != nil {
			t.Fatalf("Failed to save bookmark with category %s: %v", category, err)
		}
	}

	// Test that we can search by each custom category
	for _, category := range customCategories {
		bookmarks, err := bookmark.SearchByCategory(appInstance.Service, bookmark.CategoryType(category))
		if err != nil {
			t.Fatalf("Failed to search by category %s: %v", category, err)
		}

		if len(bookmarks) != 1 {
			t.Errorf("Expected 1 bookmark for category %s, got %d", category, len(bookmarks))
		}

		if string(bookmarks[0].Category) != category {
			t.Errorf("Expected category %s, got %s", category, bookmarks[0].Category)
		}
	}

	// Test that total count is correct
	allBookmarks, err := bookmark.List(appInstance.Service, 0, 0)
	if err != nil {
		t.Fatalf("Failed to list all bookmarks: %v", err)
	}

	if len(allBookmarks) != len(customCategories) {
		t.Errorf("Expected %d bookmarks total, got %d", len(customCategories), len(allBookmarks))
	}
}
