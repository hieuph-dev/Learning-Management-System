package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"index;size:100;not null" json:"email"`
	Token     string         `gorm:"index;size:255;not null" json:"token"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	Used      bool           `gorm:"default:false" json:"used"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Table name
func (PasswordReset) TableName() string {
	return "password_resets"
}
