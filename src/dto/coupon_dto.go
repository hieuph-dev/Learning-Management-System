package dto

import "time"

// ============ PUBLIC COUPON DTOs ============

type CheckCouponRequest struct {
	CouponCode string  `json:"coupon_code" binding:"required"`
	CourseId   uint    `json:"course_id" binding:"required"`
	OrderTotal float64 `json:"order_total" binding:"required,gt=0"`
}

type CheckCouponResponse struct {
	Valid             bool     `json:"valid"`
	CouponCode        string   `json:"coupon_code,omitempty"`
	DiscountType      string   `json:"discount_type,omitempty"`
	DiscountValue     float64  `json:"discount_value,omitempty"`
	DiscountAmount    float64  `json:"discount_amount,omitempty"`
	FinalPrice        float64  `json:"final_price,omitempty"`
	Message           string   `json:"message"`
	MinOrderAmount    float64  `json:"min_order_amount,omitempty"`
	MaxDiscountAmount *float64 `json:"max_discount_amount,omitempty"`
}

// ============ ADMIN COUPON DTOs ============

type GetAdminCouponsQueryRequest struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	IsActive   *bool  `form:"is_active"`
	SearchCode string `form:"search_code"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type AdminCouponItem struct {
	Id                uint       `json:"id"`
	Code              string     `json:"code"`
	Description       string     `json:"description"`
	DiscountType      string     `json:"discount_type"`
	DiscountValue     float64    `json:"discount_value"`
	MinOrderAmount    float64    `json:"min_order_amount"`
	MaxDiscountAmount *float64   `json:"max_discount_amount,omitempty"`
	UsageLimit        *int       `json:"usage_limit,omitempty"`
	UsedCount         int        `json:"used_count"`
	ValidFrom         *time.Time `json:"valid_from,omitempty"`
	ValidTo           *time.Time `json:"valid_to,omitempty"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
}

type GetAdminCouponsResponse struct {
	Coupons    []AdminCouponItem `json:"coupons"`
	Pagination PaginationInfo    `json:"pagination"`
}

type CreateCouponRequest struct {
	Code              string     `json:"code" binding:"required,min=3,max=50"`
	Description       string     `json:"description" binding:"max=200"`
	DiscountType      string     `json:"discount_type" binding:"required,oneof=percentage fixed"`
	DiscountValue     float64    `json:"discount_value" binding:"required,gt=0"`
	MinOrderAmount    float64    `json:"min_order_amount" binding:"omitempty,gte=0"`
	MaxDiscountAmount *float64   `json:"max_discount_amount" binding:"omitempty,gt=0"`
	UsageLimit        *int       `json:"usage_limit" binding:"omitempty,gt=0"`
	ValidFrom         *time.Time `json:"valid_from"`
	ValidTo           *time.Time `json:"valid_to"`
	IsActive          *bool      `json:"is_active"`
}

type CreateCouponResponse struct {
	Id                uint       `json:"id"`
	Code              string     `json:"code"`
	Description       string     `json:"description"`
	DiscountType      string     `json:"discount_type"`
	DiscountValue     float64    `json:"discount_value"`
	MinOrderAmount    float64    `json:"min_order_amount"`
	MaxDiscountAmount *float64   `json:"max_discount_amount,omitempty"`
	UsageLimit        *int       `json:"usage_limit,omitempty"`
	ValidFrom         *time.Time `json:"valid_from,omitempty"`
	ValidTo           *time.Time `json:"valid_to,omitempty"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	Message           string     `json:"message"`
}

type UpdateCouponRequest struct {
	Description       *string    `json:"description" binding:"omitempty,max=200"`
	DiscountType      *string    `json:"discount_type" binding:"omitempty,oneof=percentage fixed"`
	DiscountValue     *float64   `json:"discount_value" binding:"omitempty,gt=0"`
	MinOrderAmount    *float64   `json:"min_order_amount" binding:"omitempty,gte=0"`
	MaxDiscountAmount *float64   `json:"max_discount_amount" binding:"omitempty,gt=0"`
	UsageLimit        *int       `json:"usage_limit" binding:"omitempty,gt=0"`
	ValidFrom         *time.Time `json:"valid_from"`
	ValidTo           *time.Time `json:"valid_to"`
	IsActive          *bool      `json:"is_active"`
}

type UpdateCouponResponse struct {
	Id                uint       `json:"id"`
	Code              string     `json:"code"`
	Description       string     `json:"description"`
	DiscountType      string     `json:"discount_type"`
	DiscountValue     float64    `json:"discount_value"`
	MinOrderAmount    float64    `json:"min_order_amount"`
	MaxDiscountAmount *float64   `json:"max_discount_amount,omitempty"`
	UsageLimit        *int       `json:"usage_limit,omitempty"`
	UsedCount         int        `json:"used_count"`
	ValidFrom         *time.Time `json:"valid_from,omitempty"`
	ValidTo           *time.Time `json:"valid_to,omitempty"`
	IsActive          bool       `json:"is_active"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Message           string     `json:"message"`
}

type DeleteCouponResponse struct {
	Message string `json:"message"`
}
