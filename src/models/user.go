package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------------- Users ----------------
type User struct {
	Id            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email         string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password      string         `gorm:"size:255;not null" json:"-"`
	FullName      string         `gorm:"size:100;not null" json:"full_name"`
	AvatarURL     string         `gorm:"size:255" json:"avatar_url"`
	Phone         string         `gorm:"size:20" json:"phone"`
	Bio           string         `json:"bio"`
	Role          string         `gorm:"size:20;default:student" json:"role"` // admin,
	Status        string         `gorm:"size:20;default:active" json:"status"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
