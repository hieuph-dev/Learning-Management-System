package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------------- Lessons ----------------
type Lesson struct {
	Id            uint           `gorm:"primaryKey" json:"id"`
	CourseId      uint           `json:"course_id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Slug          string         `gorm:"size:200;not null" json:"slug"`
	Description   string         `json:"description"`
	Content       string         `json:"content"`
	VideoURL      string         `gorm:"size:255" json:"video_url"`
	VideoDuration int            `json:"video_duration"`
	LessonOrder   int            `gorm:"not null" json:"lesson_order"`
	IsPreview     bool           `gorm:"default:false" json:"is_preview"`
	IsPublished   bool           `gorm:"default:true" json:"is_published"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
