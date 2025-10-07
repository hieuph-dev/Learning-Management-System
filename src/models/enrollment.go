package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------------- Enrollments ----------------
type Enrollment struct {
	Id                 uint           `gorm:"primaryKey" json:"id"`
	UserId             uint           `json:"user_id"`
	User               User           `gorm:"foreignKey:UserId" json:"user"`
	CourseId           uint           `json:"course_id"`
	Course             Course         `gorm:"foreignKey:CourseId" json:"course"` // ✅ Thêm relation để Preload
	EnrolledAt         time.Time      `json:"enrolled_at"`
	CompletedAt        *time.Time     `json:"completed_at"`
	ProgressPercentage float64        `gorm:"default:0" json:"progress_percentage"`
	LastAccessedAt     *time.Time     `json:"last_accessed_at"`
	Status             string         `gorm:"size:20;default:active" json:"status"` // active, completed, dropped
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
