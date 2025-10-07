package dto

import "time"

// ============ ORDER DTOs ============

// Request tạo order
type CreateOrderRequest struct {
	CourseId   uint   `json:"course_id" binding:"required"`
	CouponCode string `json:"coupon_code"`
}

type CreateOrderResponse struct {
	OrderId        uint      `json:"order_id"`
	OrderCode      string    `json:"order_code"`
	CourseId       uint      `json:"course_id"`
	CourseTitle    string    `json:"course_title"`
	OriginalPrice  float64   `json:"original_price"`
	DiscountAmount float64   `json:"discount_amount"`
	FinalPrice     float64   `json:"final_price"`
	CouponCode     string    `json:"coupon_code,omitempty"`
	PaymentStatus  string    `json:"payment_status"`
	CreatedAt      time.Time `json:"created_at"`
	Message        string    `json:"message"`
}

// Query request cho order history
type GetOrderHistoryQueryRequest struct {
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
	PaymentStatus string `form:"payment_status" binding:"omitempty,oneof=pending paid failed cancelled"`
	SortBy        string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type OrderHistoryItem struct {
	Id              uint       `json:"id"`
	OrderCode       string     `json:"order_code"`
	CourseId        uint       `json:"course_id"`
	CourseTitle     string     `json:"course_title"`
	CourseThumbnail string     `json:"course_thumbnail"`
	OriginalPrice   float64    `json:"original_price"`
	DiscountAmount  float64    `json:"discount_amount"`
	FinalPrice      float64    `json:"final_price"`
	PaymentStatus   string     `json:"payment_status"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type GetOrderHistoryResponse struct {
	Orders     []OrderHistoryItem `json:"orders"`
	Pagination PaginationInfo     `json:"pagination"`
}

// Response lấy order detail
type OrderDetailResponse struct {
	Id              uint       `json:"id"`
	OrderCode       string     `json:"order_code"`
	UserId          uint       `json:"user_id"`
	CourseId        uint       `json:"course_id"`
	CourseTitle     string     `json:"course_title"`
	CourseThumbnail string     `json:"course_thumbnail"`
	InstructorName  string     `json:"instructor_name"`
	OriginalPrice   float64    `json:"original_price"`
	DiscountAmount  float64    `json:"discount_amount"`
	FinalPrice      float64    `json:"final_price"`
	CouponCode      string     `json:"coupon_code,omitempty"`
	PaymentStatus   string     `json:"payment_status"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Request thanh toán order
type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=credit_card paypal momo zalopay bank_transfer"`
}

type PayOrderResponse struct {
	OrderId       uint      `json:"order_id"`
	OrderCode     string    `json:"order_code"`
	PaymentStatus string    `json:"payment_status"`
	PaymentMethod string    `json:"payment_method"`
	PaidAt        time.Time `json:"paid_at"`
	Message       string    `json:"message"`
}

// Request validate coupon
type ValidateCouponRequest struct {
	CouponCode string  `json:"coupon_code" binding:"required"`
	CourseId   uint    `json:"course_id" binding:"required"`
	OrderTotal float64 `json:"order_total" binding:"required,gt=0"`
}

type ValidateCouponResponse struct {
	Valid             bool     `json:"valid"`
	CouponCode        string   `json:"coupon_code"`
	DiscountType      string   `json:"discount_type"`
	DiscountValue     float64  `json:"discount_value"`
	DiscountAmount    float64  `json:"discount_amount"`
	FinalPrice        float64  `json:"final_price"`
	Message           string   `json:"message"`
	MinOrderAmount    float64  `json:"min_order_amount,omitempty"`
	MaxDiscountAmount *float64 `json:"max_discount_amount,omitempty"`
}

// ============ ADMIN ORDER DTOs ============

// Query request cho admin orders
type GetAdminOrdersQueryRequest struct {
	Page          int    `form:"page" binding:"omitempty,min=1"`
	Limit         int    `form:"limit" binding:"omitempty,min=1,max=100"`
	UserId        *uint  `form:"user_id" binding:"omitempty"`
	CourseId      *uint  `form:"course_id" binding:"omitempty"`
	PaymentStatus string `form:"payment_status" binding:"omitempty,oneof=pending paid failed cancelled"`
	PaymentMethod string `form:"payment_method" binding:"omitempty,oneof=credit_card paypal momo zalopay bank_transfer"`
	Search        string `form:"search" binding:"omitempty,search"`
	OrderBy       string `form:"order_by" binding:"omitempty,oneof=created_at updated_at final_price"`
	SortBy        string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
	DateFrom      string `form:"date_from" binding:"omitempty,datetime=2006-01-02"`
	DateTo        string `form:"date_to" binding:"omitempty,datetime=2006-01-02"`
}

type AdminOrderItem struct {
	Id              uint       `json:"id"`
	OrderCode       string     `json:"order_code"`
	UserId          uint       `json:"user_id"`
	Username        string     `json:"username"`
	UserEmail       string     `json:"user_email"`
	CourseId        uint       `json:"course_id"`
	CourseTitle     string     `json:"course_title"`
	CourseThumbnail string     `json:"course_thumbnail"`
	InstructorName  string     `json:"instructor_name"`
	OriginalPrice   float64    `json:"original_price"`
	DiscountAmount  float64    `json:"discount_amount"`
	FinalPrice      float64    `json:"final_price"`
	CouponCode      string     `json:"coupon_code,omitempty"`
	PaymentMethod   string     `json:"payment_method"`
	PaymentStatus   string     `json:"payment_status"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type GetAdminOrdersResponse struct {
	Orders     []AdminOrderItem `json:"orders"`
	Pagination PaginationInfo   `json:"pagination"`
	Statistics OrderStatistics  `json:"statistics"`
}

type OrderStatistics struct {
	TotalOrders       int     `json:"total_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingOrders     int     `json:"pending_orders"`
	CompletedOrders   int     `json:"completed_orders"`
	FailedOrders      int     `json:"failed_orders"`
	CancelledOrders   int     `json:"cancelled_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// Request update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending paid failed cancelled refunded"`
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

type UpdateOrderStatusResponse struct {
	Id            uint      `json:"id"`
	OrderCode     string    `json:"order_code"`
	PaymentStatus string    `json:"payment_status"`
	UpdatedAt     time.Time `json:"updated_at"`
	Message       string    `json:"message"`
}
