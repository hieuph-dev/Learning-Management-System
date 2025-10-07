package dto

import "time"

// EnrollCourseRequest - Request để enroll course
type EnrollCourseRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=credit_card paypal momo zalopay bank_transfer"`
	CouponCode    string `json:"coupon_code" binding:"omitempty"`
}

// EnrollCourseResponse - Response sau khi enroll thành công
type EnrollCourseResponse struct {
	EnrollmentId   uint      `json:"enrollment_id"`
	OrderId        uint      `json:"order_id"`
	OrderCode      string    `json:"order_code"`
	CourseId       uint      `json:"course_id"`
	CourseTitle    string    `json:"course_title"`
	OriginalPrice  float64   `json:"original_price"`
	DiscountAmount float64   `json:"discount_amount"`
	FinalPrice     float64   `json:"final_price"`
	PaymentMethod  string    `json:"payment_method"`
	PaymentStatus  string    `json:"payment_status"`
	EnrolledAt     time.Time `json:"enrolled_at"`
	Message        string    `json:"message"`
}

// CheckEnrollmentResponse - Kiểm tra user đã enroll chưa
type CheckEnrollmentResponse struct {
	IsEnrolled      bool       `json:"is_enrolled"`
	EnrollmentId    uint       `json:"enrollment_id,omitempty"`
	EnrolledAt      *time.Time `json:"enrolled_at,omitempty"`
	ProgressPercent float64    `json:"progress_percent,omitempty"`
}

// GetMyEnrollmentsQueryRequest - Query parameters cho danh sách enrollments
type GetMyEnrollmentsQueryRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=50"`
	Status string `form:"status" binding:"omitempty,oneof=active completed dropped"`
}

// EnrollmentItem - Thông tin enrollment
type EnrollmentItem struct {
	Id                 uint      `json:"id"`
	CourseId           uint      `json:"course_id"`
	CourseTitle        string    `json:"course_title"`
	CourseThumbnail    string    `json:"course_thumbnail"`
	InstructorName     string    `json:"instructor_name"`
	EnrolledAt         time.Time `json:"enrolled_at"`
	ProgressPercentage float64   `json:"progress_percentage"`
	Status             string    `json:"status"`
	TotalLessons       int       `json:"total_lessons"`
	CompletedLessons   int       `json:"completed_lessons"`
}

// GetMyEnrollmentsResponse - Response danh sách enrollments
type GetMyEnrollmentsResponse struct {
	Enrollments []EnrollmentItem `json:"enrollments"`
	Pagination  PaginationInfo   `json:"pagination"`
}
