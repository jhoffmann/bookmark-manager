package bookmark

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jhoffmann/bookmark-manager/internal/config"
	"github.com/jhoffmann/bookmark-manager/internal/database"
)

// setupTestDatabase creates a test database for testing
func setupTestDatabase(t *testing.T) (database.DB, func()) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_bookmarks.db")

	cfg := &config.Config{
		DatabasePath: dbPath,
		LogLevel:     "silent",
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.RemoveAll(tempDir)
	}

	return db, cleanup
}

func TestNewService(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)
	if service == nil {
		t.Error("NewService() returned nil")
	}

	if service.db != db {
		t.Error("NewService() did not set database correctly")
	}
}

func TestBookmark_BeforeCreate(t *testing.T) {
	tests := []struct {
		name     string
		bookmark *Bookmark
		wantCat  CategoryType
	}{
		{
			name:     "preserves empty category",
			bookmark: &Bookmark{Folder: "/test"},
			wantCat:  "",
		},
		{
			name:     "preserves existing category",
			bookmark: &Bookmark{Folder: "/test", Category: "work"},
			wantCat:  "work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalTime := tt.bookmark.DateCreated

			err := tt.bookmark.BeforeCreate(nil)
			if err != nil {
				t.Errorf("BeforeCreate() error = %v", err)
			}

			if tt.bookmark.Category != tt.wantCat {
				t.Errorf("BeforeCreate() category = %v, want %v", tt.bookmark.Category, tt.wantCat)
			}

			if originalTime.IsZero() && tt.bookmark.DateCreated.IsZero() {
				t.Error("BeforeCreate() should set DateCreated when it's zero")
			}
		})
	}
}

func TestBookmark_Save_Create(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)
	bookmark := &Bookmark{
		Folder:   "/home/user/documents",
		Category: "personal",
	}

	err := bookmark.Save(service)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if bookmark.ID == 0 {
		t.Error("Save() should set ID for new bookmark")
	}

	if bookmark.DateCreated.IsZero() {
		t.Error("Save() should set DateCreated")
	}
}

func TestBookmark_Save_Update(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	// Create initial bookmark
	bookmark := &Bookmark{
		Folder:   "/home/user/documents",
		Category: "personal",
	}

	err := bookmark.Save(service)
	if err != nil {
		t.Fatalf("Initial Save() error = %v", err)
	}

	originalID := bookmark.ID

	// Update the bookmark
	bookmark.Category = "work"
	bookmark.Folder = "/home/user/work"

	err = bookmark.Save(service)
	if err != nil {
		t.Fatalf("Update Save() error = %v", err)
	}

	if bookmark.ID != originalID {
		t.Error("Save() should preserve ID during update")
	}

	// Verify update in database
	retrieved, err := GetByID(service, bookmark.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if retrieved.Category != "work" {
		t.Errorf("Update not persisted: category = %v, want %v", retrieved.Category, "work")
	}

	if retrieved.Folder != "/home/user/work" {
		t.Errorf("Update not persisted: folder = %v, want %v", retrieved.Folder, "/home/user/work")
	}
}

func TestBookmark_Save_Validation(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	tests := []struct {
		name     string
		bookmark *Bookmark
		wantErr  bool
	}{
		{
			name:     "empty folder",
			bookmark: &Bookmark{Category: "work"},
			wantErr:  true,
		},
		{
			name:     "custom category allowed",
			bookmark: &Bookmark{Folder: "/test", Category: "my-custom-category"},
			wantErr:  false,
		},
		{
			name:     "valid bookmark with category",
			bookmark: &Bookmark{Folder: "/test", Category: "work"},
			wantErr:  false,
		},
		{
			name:     "valid bookmark with empty category",
			bookmark: &Bookmark{Folder: "/test", Category: ""},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bookmark.Save(service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check that empty category remains empty
			if tt.name == "valid bookmark with empty category" && tt.bookmark.Category != "" {
				t.Errorf("Save() should preserve empty category, got %v", tt.bookmark.Category)
			}
		})
	}
}

func TestBookmark_Delete(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	// Create a bookmark to delete
	bookmark := &Bookmark{
		Folder:   "/home/user/temp",
		Category: "work",
	}

	err := bookmark.Save(service)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	bookmarkID := bookmark.ID

	// Delete the bookmark
	err = bookmark.Delete(service)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it's deleted (should not be found)
	_, err = GetByID(service, bookmarkID)
	if err == nil {
		t.Error("GetByID() should return error for deleted bookmark")
	}
}

func TestBookmark_Delete_NoID(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)
	bookmark := &Bookmark{Folder: "/test"}

	err := bookmark.Delete(service)
	if err == nil {
		t.Error("Delete() should return error when ID is not set")
	}
}

func TestSearchByCategory(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	// Create test bookmarks with different categories
	bookmarks := []*Bookmark{
		{Folder: "/work/project1", Category: "work"},
		{Folder: "/work/project2", Category: "work"},
		{Folder: "/personal/docs", Category: "personal"},
		{Folder: "/default/folder", Category: "misc"},
	}

	for _, b := range bookmarks {
		if err := b.Save(service); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	// Test searching by work category
	workBookmarks, err := SearchByCategory(service, "work")
	if err != nil {
		t.Fatalf("SearchByCategory() error = %v", err)
	}

	if len(workBookmarks) != 2 {
		t.Errorf("SearchByCategory(work) = %d bookmarks, want 2", len(workBookmarks))
	}

	// Test searching by personal category
	personalBookmarks, err := SearchByCategory(service, "personal")
	if err != nil {
		t.Fatalf("SearchByCategory() error = %v", err)
	}

	if len(personalBookmarks) != 1 {
		t.Errorf("SearchByCategory(personal) = %d bookmarks, want 1", len(personalBookmarks))
	}

	// Test method on bookmark instance
	testBookmark := &Bookmark{}
	instanceResults, err := testBookmark.SearchByCategory(service, "work")
	if err != nil {
		t.Fatalf("Bookmark.SearchByCategory() error = %v", err)
	}

	if len(instanceResults) != len(workBookmarks) {
		t.Errorf("Instance method returned different results than static method")
	}
}

func TestGetByID(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	// Create a bookmark
	original := &Bookmark{
		Folder:   "/home/user/photos",
		Category: "personal",
	}

	err := original.Save(service)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Retrieve by ID
	retrieved, err := GetByID(service, original.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if retrieved.ID != original.ID {
		t.Errorf("GetByID() ID = %d, want %d", retrieved.ID, original.ID)
	}

	if retrieved.Folder != original.Folder {
		t.Errorf("GetByID() Folder = %s, want %s", retrieved.Folder, original.Folder)
	}

	if retrieved.Category != original.Category {
		t.Errorf("GetByID() Category = %s, want %s", retrieved.Category, original.Category)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	_, err := GetByID(service, 999)
	if err == nil {
		t.Error("GetByID() should return error for non-existent ID")
	}
}

func TestList(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	// Create test bookmarks
	bookmarks := []*Bookmark{
		{Folder: "/folder1", Category: "work"},
		{Folder: "/folder2", Category: "work"},
		{Folder: "/folder3", Category: "personal"},
	}

	for _, b := range bookmarks {
		if err := b.Save(service); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
		// Small delay to ensure different DateCreated times
		time.Sleep(time.Millisecond)
	}

	// Test listing all bookmarks
	allBookmarks, err := List(service, 0, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(allBookmarks) != 3 {
		t.Errorf("List() = %d bookmarks, want 3", len(allBookmarks))
	}

	// Test with limit
	limitedBookmarks, err := List(service, 2, 0)
	if err != nil {
		t.Fatalf("List() with limit error = %v", err)
	}

	if len(limitedBookmarks) != 2 {
		t.Errorf("List() with limit = %d bookmarks, want 2", len(limitedBookmarks))
	}

	// Test with offset
	offsetBookmarks, err := List(service, 0, 1)
	if err != nil {
		t.Fatalf("List() with offset error = %v", err)
	}

	if len(offsetBookmarks) != 2 {
		t.Errorf("List() with offset = %d bookmarks, want 2", len(offsetBookmarks))
	}
}

func TestSearchByFolder(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	// Create test bookmarks with different folder paths
	bookmarks := []*Bookmark{
		{Folder: "/home/user/documents", Category: "personal"},
		{Folder: "/home/user/downloads", Category: "personal"},
		{Folder: "/work/projects", Category: "work"},
		{Folder: "/tmp/temp", Category: "work"},
	}

	for _, b := range bookmarks {
		if err := b.Save(service); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	// Test searching by partial folder path
	homeBookmarks, err := SearchByFolder(service, "home/user")
	if err != nil {
		t.Fatalf("SearchByFolder() error = %v", err)
	}

	if len(homeBookmarks) != 2 {
		t.Errorf("SearchByFolder('home/user') = %d bookmarks, want 2", len(homeBookmarks))
	}

	// Test searching by specific folder
	workBookmarks, err := SearchByFolder(service, "work")
	if err != nil {
		t.Fatalf("SearchByFolder() error = %v", err)
	}

	if len(workBookmarks) != 1 {
		t.Errorf("SearchByFolder('work') = %d bookmarks, want 1", len(workBookmarks))
	}

	// Test searching with no matches
	noMatches, err := SearchByFolder(service, "nonexistent")
	if err != nil {
		t.Fatalf("SearchByFolder() error = %v", err)
	}

	if len(noMatches) != 0 {
		t.Errorf("SearchByFolder('nonexistent') = %d bookmarks, want 0", len(noMatches))
	}
}

func TestBookmark_Validate(t *testing.T) {
	tests := []struct {
		name     string
		bookmark *Bookmark
		wantErr  bool
	}{
		{
			name:     "valid bookmark",
			bookmark: &Bookmark{Folder: "/test", Category: "work"},
			wantErr:  false,
		},
		{
			name:     "empty folder",
			bookmark: &Bookmark{Category: "work"},
			wantErr:  true,
		},
		{
			name:     "custom category allowed",
			bookmark: &Bookmark{Folder: "/test", Category: "my-custom-category"},
			wantErr:  false,
		},
		{
			name:     "empty category allowed",
			bookmark: &Bookmark{Folder: "/test", Category: ""},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bookmark.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check that empty category remains empty
			if tt.name == "empty category allowed" && tt.bookmark.Category != "" {
				t.Errorf("validate() should preserve empty category, got %v", tt.bookmark.Category)
			}
		})
	}
}

func TestBookmark_String(t *testing.T) {
	bookmark := &Bookmark{
		ID:          1,
		Folder:      "/test/folder",
		Category:    "work",
		DateCreated: time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC),
	}

	str := bookmark.String()
	expected := "Bookmark{ID: 1, Folder: /test/folder, Category: work, DateCreated: 2023-12-25 15:30:00}"

	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}

func TestDefaultCategory(t *testing.T) {
	// Test that CategoryType can be empty
	var emptyCategory CategoryType = ""
	if string(emptyCategory) != "" {
		t.Errorf("Empty CategoryType should be empty string, got %v", emptyCategory)
	}
}

// Integration test that verifies end-to-end functionality
func TestBookmark_Integration(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	service := NewService(db)

	// Create multiple bookmarks with user-defined categories
	bookmarks := []*Bookmark{
		{Folder: "/work/project1", Category: "work"},
		{Folder: "/home/documents", Category: "personal"},
		{Folder: "/tmp", Category: "work"},
		{Folder: "/study/materials", Category: "education"},
	}

	// Save all bookmarks
	for _, b := range bookmarks {
		if err := b.Save(service); err != nil {
			t.Fatalf("Failed to save bookmark: %v", err)
		}
		if b.ID == 0 {
			t.Error("Bookmark ID should be set after save")
		}
	}

	// Test listing
	allBookmarks, err := List(service, 0, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(allBookmarks) != 4 {
		t.Errorf("Expected 4 bookmarks, got %d", len(allBookmarks))
	}

	// Test search by category
	workBookmarks, err := SearchByCategory(service, "work")
	if err != nil {
		t.Fatalf("SearchByCategory failed: %v", err)
	}
	if len(workBookmarks) != 2 {
		t.Errorf("Expected 2 work bookmarks, got %d", len(workBookmarks))
	}

	// Test search by custom category
	eduBookmarks, err := SearchByCategory(service, "education")
	if err != nil {
		t.Fatalf("SearchByCategory failed: %v", err)
	}
	if len(eduBookmarks) != 1 {
		t.Errorf("Expected 1 education bookmark, got %d", len(eduBookmarks))
	}

	// Test search by folder
	homeBookmarks, err := SearchByFolder(service, "home")
	if err != nil {
		t.Fatalf("SearchByFolder failed: %v", err)
	}
	if len(homeBookmarks) != 1 {
		t.Errorf("Expected 1 bookmark with 'home' in path, got %d", len(homeBookmarks))
	}

	// Test update with new custom category
	workBookmarks[0].Folder = "/work/updated-project"
	workBookmarks[0].Category = "work-priority"
	if err := workBookmarks[0].Save(service); err != nil {
		t.Fatalf("Failed to update bookmark: %v", err)
	}

	// Verify update
	updated, err := GetByID(service, workBookmarks[0].ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated bookmark: %v", err)
	}
	if updated.Folder != "/work/updated-project" {
		t.Errorf("Update failed: folder = %v, want '/work/updated-project'", updated.Folder)
	}
	if updated.Category != "work-priority" {
		t.Errorf("Update failed: category = %v, want 'work-priority'", updated.Category)
	}

	// Test delete
	if err := updated.Delete(service); err != nil {
		t.Fatalf("Failed to delete bookmark: %v", err)
	}

	// Verify deletion
	_, err = GetByID(service, updated.ID)
	if err == nil {
		t.Error("Deleted bookmark should not be found")
	}

	// Final count should be 3
	finalBookmarks, err := List(service, 0, 0)
	if err != nil {
		t.Fatalf("Final list failed: %v", err)
	}
	if len(finalBookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks after delete, got %d", len(finalBookmarks))
	}
}
