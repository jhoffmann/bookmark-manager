// Package bookmark provides functionality for managing folder bookmarks.
// It handles CRUD operations for bookmarks that track folders with categories.
package bookmark

import (
	"fmt"
	"time"

	"github.com/jhoffmann/bookmark-manager/internal/database"
	"gorm.io/gorm"
)

// CategoryType represents the category of a bookmark
type CategoryType string

// Bookmark represents a folder bookmark entry
type Bookmark struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Folder      string         `gorm:"not null" json:"folder"`
	DateCreated time.Time      `json:"date_created"`
	Category    CategoryType   `gorm:"type:varchar(50)" json:"category"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Service provides bookmark operations with database access
type Service struct {
	db database.DB
}

// NewService creates a new bookmark service with the provided database
func NewService(db database.DB) *Service {
	return &Service{db: db}
}

// BeforeCreate is a GORM hook that runs before creating a bookmark
func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.DateCreated.IsZero() {
		b.DateCreated = time.Now()
	}
	// Allow empty category - no default assignment
	return nil
}

// Save saves the bookmark to the database
func (b *Bookmark) Save(service *Service) error {
	if err := b.validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	gormDB := service.db.GetDB()
	if gormDB == nil {
		return fmt.Errorf("database connection is not available")
	}

	if b.ID == 0 {
		// Create new bookmark
		if err := gormDB.Create(b).Error; err != nil {
			return fmt.Errorf("failed to create bookmark: %w", err)
		}
	} else {
		// Update existing bookmark
		if err := gormDB.Save(b).Error; err != nil {
			return fmt.Errorf("failed to update bookmark: %w", err)
		}
	}

	return nil
}

// Delete removes the bookmark from the database (soft delete)
func (b *Bookmark) Delete(service *Service) error {
	if b.ID == 0 {
		return fmt.Errorf("cannot delete bookmark: ID is required")
	}

	gormDB := service.db.GetDB()
	if gormDB == nil {
		return fmt.Errorf("database connection is not available")
	}

	if err := gormDB.Delete(b).Error; err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	return nil
}

// SearchByCategory searches for bookmarks by category
func (b *Bookmark) SearchByCategory(service *Service, category CategoryType) ([]*Bookmark, error) {
	return SearchByCategory(service, category)
}

// SearchByCategory searches for bookmarks by category (static method)
func SearchByCategory(service *Service, category CategoryType) ([]*Bookmark, error) {
	gormDB := service.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmarks []*Bookmark
	if err := gormDB.Where("category = ?", category).Order("date_created DESC").Find(&bookmarks).Error; err != nil {
		return nil, fmt.Errorf("failed to search bookmarks by category: %w", err)
	}

	return bookmarks, nil
}

// GetByID retrieves a bookmark by its ID
func GetByID(service *Service, id uint) (*Bookmark, error) {
	gormDB := service.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmark Bookmark
	if err := gormDB.First(&bookmark, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bookmark with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get bookmark: %w", err)
	}

	return &bookmark, nil
}

// List retrieves all bookmarks with optional limit and offset
func List(service *Service, limit, offset int) ([]*Bookmark, error) {
	gormDB := service.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmarks []*Bookmark
	query := gormDB.Order("date_created DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&bookmarks).Error; err != nil {
		return nil, fmt.Errorf("failed to list bookmarks: %w", err)
	}

	return bookmarks, nil
}

// SearchByFolder searches for bookmarks by folder path (partial match)
func SearchByFolder(service *Service, folderPath string) ([]*Bookmark, error) {
	gormDB := service.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmarks []*Bookmark
	if err := gormDB.Where("folder LIKE ?", "%"+folderPath+"%").Order("date_created DESC").Find(&bookmarks).Error; err != nil {
		return nil, fmt.Errorf("failed to search bookmarks by folder: %w", err)
	}

	return bookmarks, nil
}

// validate performs validation on the bookmark fields
func (b *Bookmark) validate() error {
	if b.Folder == "" {
		return fmt.Errorf("folder path is required")
	}

	// Allow empty category - no default assignment

	return nil
}

// String returns a string representation of the bookmark
func (b *Bookmark) String() string {
	return fmt.Sprintf("Bookmark{ID: %d, Folder: %s, Category: %s, DateCreated: %s}",
		b.ID, b.Folder, b.Category, b.DateCreated.Format("2006-01-02 15:04:05"))
}
