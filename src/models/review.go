package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	Id          uint           `gorm:"primaryKey" json:"id"`
	UserId      uint           `json:"user_id"`
	User        User           `gorm:"foreignKey:UserId" json:"user"`
	CourseId    uint           `json:"course_id"`
	Course      Course         `gorm:"foreignKey:CourseId" json:"course"`
	Rating      int            `gorm:"not null" json:"rating"`
	Comment     string         `json:"comment"`
	IsPublished bool           `gorm:"default:true" json:"is_published"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
