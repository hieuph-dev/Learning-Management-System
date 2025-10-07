package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------------- Coupons ----------------
type Coupon struct {
	Id                uint           `gorm:"primaryKey" json:"id"`
	Code              string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description       string         `gorm:"size:200" json:"description"`
	DiscountType      string         `gorm:"size:20" json:"discount_type"` // percentage, fixed
	DiscountValue     float64        `gorm:"not null" json:"discount_value"`
	MinOrderAmount    float64        `gorm:"default:0" json:"min_order_amount"`
	MaxDiscountAmount *float64       `json:"max_discount_amount"`
	UsageLimit        *int           `json:"usage_limit"`
	UsedCount         int            `gorm:"default:0" json:"used_count"`
	ValidFrom         *time.Time     `json:"valid_from"`
	ValidTo           *time.Time     `json:"valid_to"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}
