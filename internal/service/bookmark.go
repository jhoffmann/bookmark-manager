// Package service provides business logic services for the bookmark manager application.
package service

import (
	"fmt"

	"github.com/jhoffmann/bookmark-manager/internal/database"
	"github.com/jhoffmann/bookmark-manager/internal/models"
	"gorm.io/gorm"
)

// Bookmarks provides bookmark operations with database access
type Bookmarks struct {
	db database.DB
}

// NewBookmarks creates a new bookmark service with the provided database
func NewBookmarks(db database.DB) *Bookmarks {
	return &Bookmarks{db: db}
}

// Save saves the bookmark to the database
func (s *Bookmarks) Save(b *models.Bookmark) error {
	if err := b.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	gormDB := s.db.GetDB()
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
func (s *Bookmarks) Delete(b *models.Bookmark) error {
	if b.ID == 0 {
		return fmt.Errorf("cannot delete bookmark: ID is required")
	}

	gormDB := s.db.GetDB()
	if gormDB == nil {
		return fmt.Errorf("database connection is not available")
	}

	if err := gormDB.Delete(b).Error; err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	return nil
}

// GetByID retrieves a bookmark by its ID
func (s *Bookmarks) GetByID(id uint) (*models.Bookmark, error) {
	gormDB := s.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmark models.Bookmark
	if err := gormDB.First(&bookmark, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bookmark with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get bookmark: %w", err)
	}

	return &bookmark, nil
}

// List retrieves all bookmarks with optional limit and offset
func (s *Bookmarks) List(limit, offset int) ([]*models.Bookmark, error) {
	gormDB := s.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmarks []*models.Bookmark
	query := gormDB.Order("category, folder")

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

// SearchByCategory searches for bookmarks by category
func (s *Bookmarks) SearchByCategory(category models.CategoryType) ([]*models.Bookmark, error) {
	gormDB := s.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmarks []*models.Bookmark
	if err := gormDB.Where("category = ?", category).Order("category, folder").Find(&bookmarks).Error; err != nil {
		return nil, fmt.Errorf("failed to search bookmarks by category: %w", err)
	}

	return bookmarks, nil
}

// SearchByFolder searches for bookmarks by folder path (partial match)
func (s *Bookmarks) SearchByFolder(folderPath string) ([]*models.Bookmark, error) {
	gormDB := s.db.GetDB()
	if gormDB == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	var bookmarks []*models.Bookmark
	if err := gormDB.Where("folder LIKE ?", "%"+folderPath+"%").Order("category, folder").Find(&bookmarks).Error; err != nil {
		return nil, fmt.Errorf("failed to search bookmarks by folder: %w", err)
	}

	return bookmarks, nil
}
