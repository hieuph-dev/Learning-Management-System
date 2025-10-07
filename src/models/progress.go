package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------------- Progress ----------------
type Progress struct {
	Id            uint           `gorm:"primaryKey" json:"id"`
	UserId        uint           `json:"user_id"`
	LessonId      uint           `json:"lesson_id"`
	CourseId      uint           `json:"course_id"`
	IsCompleted   bool           `gorm:"default:false" json:"is_completed"`
	CompletedAt   *time.Time     `json:"completed_at"`
	WatchDuration int            `gorm:"default:0" json:"watch_duration"`
	LastPosition  int            `gorm:"default:0" json:"last_position"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
