package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	Id              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"size:200;not null" json:"title"`
	Slug            string         `gorm:"uniqueIndex;size:200;not null" json:"slug"`
	Description     string         `json:"description"`
	ShortDesc       string         `gorm:"size:500" json:"short_description"`
	ThumbnailURL    string         `gorm:"size:255" json:"thumbnail_url"`
	VideoPreviewURL string         `gorm:"size:255" json:"video_preview_url"`
	Price           float64        `gorm:"not null;default:0" json:"price"`
	DiscountPrice   *float64       `json:"discount_price"`
	InstructorId    uint           `json:"instructor_id"`
	Instructor      User           `gorm:"foreignKey:InstructorId" json:"instructor"`
	CategoryId      uint           `json:"category_id"`
	Category        Category       `gorm:"foreignKey:CategoryId" json:"category"`
	Level           string         `gorm:"size:20" json:"level"` // beginner, intermediate, advanced
	DurationHours   int            `gorm:"default:0" json:"duration_hours"`
	TotalLessons    int            `gorm:"default:0" json:"total_lessons"`
	Language        string         `gorm:"size:10;default:vi" json:"language"`
	Requirements    string         `json:"requirements"`
	WhatYouLearn    string         `json:"what_you_learn"`
	Status          string         `gorm:"size:20;default:draft" json:"status"` // draft, published, archived
	IsFeatured      bool           `gorm:"default:false" json:"is_featured"`
	RatingAvg       float32        `gorm:"default:0" json:"rating_avg"`
	RatingCount     int            `gorm:"default:0" json:"rating_count"`
	EnrolledCount   int            `gorm:"default:0" json:"enrolled_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
