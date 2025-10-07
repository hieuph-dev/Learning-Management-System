package dto

import "time"

type GetUsersQueryRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Role    string `form:"role" binding:"omitempty,oneof=admin student instructor guest"`
	Status  string `form:"status" binding:"omitempty,oneof=active inactive banned"`
	Search  string `form:"search" binding:"omitempty,search"`
	OrderBy string `form:"order_by" binding:"omitempty,oneof=created_at updated_at username email"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type AdminUserItem struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GetUsersResponse struct {
	Users      []AdminUserItem `json:"users"`
	Pagination PaginationInfo  `json:"pagination"`
}

type PaginationInfo struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type AdminUserDetail struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Bio           string    `json:"bio"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	FullName      string `json:"full_name" binding:"omitempty,min=2,max=100"`
	Phone         string `json:"phone" binding:"omitempty,max=20"`
	Bio           string `json:"bio" binding:"omitempty,max=500"`
	AvatarURL     string `json:"avatar_url" binding:"omitempty,url"`
	Role          string `json:"role" binding:"omitempty,oneof=admin student instructor guest"`
	Status        string `json:"status" binding:"omitempty,oneof=active inactive banned"`
	EmailVerified bool   `json:"email_verified"`
}

type UpdateUserResponse struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Bio           string    `json:"bio"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
	UserId  uint   `json:"user_id"`
}

type ChangeUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive banned"`
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

type ChangeUserStatusResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type GetAdminCoursesQueryRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`
	Limit        int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Status       string `form:"status" binding:"omitempty,course_status"`
	Level        string `form:"level" binding:"omitempty,course_level"`
	CategoryId   uint   `form:"category_id" binding:"omitempty"`
	InstructorId uint   `form:"instructor_id" binding:"omitempty"`
	IsFeatured   *bool  `form:"is_featured" binding:"omitempty"`
	Search       string `form:"search" binding:"omitempty,search"`
	OrderBy      string `form:"order_by" binding:"omitempty,oneof=created_at updated_at title price enrolled_count rating_avg"`
	SortBy       string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type AdminCourseItem struct {
	Id             uint      `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	Price          float64   `json:"price"`
	DiscountPrice  *float64  `json:"discount_price"`
	InstructorId   uint      `json:"instructor_id"`
	InstructorName string    `json:"instructor_name"`
	CategoryId     uint      `json:"category_id"`
	CategoryName   string    `json:"category_name"`
	Level          string    `json:"level"`
	Status         string    `json:"status"`
	TotalLessons   int       `json:"total_lessons"`
	DurationHours  int       `json:"duration_hours"`
	EnrolledCount  int       `json:"enrolled_count"`
	RatingAvg      float32   `json:"rating_avg"`
	RatingCount    int       `json:"rating_count"`
	IsFeatured     bool      `json:"is_featured"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type GetAdminCoursesResponse struct {
	Courses    []AdminCourseItem `json:"courses"`
	Pagination PaginationInfo    `json:"pagination"`
}

type ChangeCourseStatusRequest struct {
	Status string `json:"status" binding:"required,course_status"`
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

type ChangeCourseStatusResponse struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
