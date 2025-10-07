package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------------- Categories ----------------
type Category struct {
	Id          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;size:100;not null" json:"slug"`
	Description string         `json:"description"`
	ImageURL    string         `gorm:"size:255" json:"image_url"`
	ParentId    *uint          `json:"parent_id"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
