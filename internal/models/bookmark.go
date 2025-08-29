// Package models provides data models for the bookmark manager application.
package models

import (
	"fmt"
	"time"

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

// BeforeCreate is a GORM hook that runs before creating a bookmark
func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.DateCreated.IsZero() {
		b.DateCreated = time.Now()
	}
	// Allow empty category - no default assignment
	return nil
}

// String returns a string representation of the bookmark
func (b *Bookmark) String() string {
	return fmt.Sprintf("Bookmark{ID: %d, Folder: %s, Category: %s, DateCreated: %s}",
		b.ID, b.Folder, b.Category, b.DateCreated.Format("2006-01-02 15:04:05"))
}

// Validate performs validation on the bookmark fields
func (b *Bookmark) Validate() error {
	if b.Folder == "" {
		return fmt.Errorf("folder path is required")
	}

	// Allow empty category - no default assignment

	return nil
}
