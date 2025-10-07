// lms/src/dto/instructor_dto.go
package dto

import "time"

// GET /api/v1/instructor/courses - Query parameters
type GetInstructorCoursesQueryRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Status  string `form:"status" binding:"omitempty,oneof=draft published archived"`
	Search  string `form:"search" binding:"omitempty,search"`
	OrderBy string `form:"order_by" binding:"omitempty,oneof=created_at updated_at title enrolled_count rating_avg"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

// Response item cho mỗi course
type InstructorCourseItem struct {
	Id            uint      `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	Price         float64   `json:"price"`
	DiscountPrice *float64  `json:"discount_price"`
	CategoryId    uint      `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	Level         string    `json:"level"`
	Status        string    `json:"status"`
	TotalLessons  int       `json:"total_lessons"`
	DurationHours int       `json:"duration_hours"`
	EnrolledCount int       `json:"enrolled_count"`
	RatingAvg     float32   `json:"rating_avg"`
	RatingCount   int       `json:"rating_count"`
	IsFeatured    bool      `json:"is_featured"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Response cho danh sách courses
type GetInstructorCoursesResponse struct {
	Courses    []InstructorCourseItem `json:"courses"`
	Pagination PaginationInfo         `json:"pagination"`
}

type CreateCourseRequest struct {
	Title         string   `json:"title" binding:"required,min=5,max=200"`
	Description   string   `json:"description" binding:"required,min=20"`
	ShortDesc     string   `json:"short_description" binding:"required,min=10,max=500"`
	CategoryId    uint     `json:"category_id" binding:"required"`
	Level         string   `json:"level" binding:"required,course_level"`
	Language      string   `json:"language" binding:"required,language_code"`
	Price         float64  `json:"price" binding:"required,positive_float"`
	DiscountPrice *float64 `json:"discount_price" binding:"omitempty,positive_float"`
	Requirements  string   `json:"requirements" binding:"omitempty"`
	WhatYouLearn  string   `json:"what_you_learn" binding:"omitempty"`
	DurationHours int      `json:"duration_hours" binding:"omitempty,min_int=0"`
}

type CreateCourseResponse struct {
	Id              uint      `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Description     string    `json:"description"`
	ShortDesc       string    `json:"short_description"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	VideoPreviewURL string    `json:"video_preview_url"`
	Price           float64   `json:"price"`
	DiscountPrice   *float64  `json:"discount_price"`
	InstructorId    uint      `json:"instructor_id"`
	CategoryId      uint      `json:"category_id"`
	CategoryName    string    `json:"category_name"`
	Level           string    `json:"level"`
	DurationHours   int       `json:"duration_hours"`
	TotalLessons    int       `json:"total_lessons"`
	Language        string    `json:"language"`
	Requirements    string    `json:"requirements"`
	WhatYouLearn    string    `json:"what_you_learn"`
	Status          string    `json:"status"`
	IsFeatured      bool      `json:"is_featured"`
	RatingAvg       float32   `json:"rating_avg"`
	RatingCount     int       `json:"rating_count"`
	EnrolledCount   int       `json:"enrolled_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpdateCourseRequest struct {
	Title         string   `json:"title" binding:"omitempty,min=5,max=200"`
	Description   string   `json:"description" binding:"omitempty,min=20"`
	ShortDesc     string   `json:"short_description" binding:"omitempty,min=10,max=500"`
	CategoryId    uint     `json:"category_id" binding:"omitempty"`
	Level         string   `json:"level" binding:"omitempty,course_level"`
	Language      string   `json:"language" binding:"omitempty,language_code"`
	Price         float64  `json:"price" binding:"omitempty,positive_float"`
	DiscountPrice *float64 `json:"discount_price" binding:"omitempty,positive_float"`
	Requirements  string   `json:"requirements"`
	WhatYouLearn  string   `json:"what_you_learn"`
	DurationHours int      `json:"duration_hours" binding:"omitempty,min_int=0"`
	Status        string   `json:"status" binding:"omitempty,course_status"`
	IsFeatured    *bool    `json:"is_featured"`
}

type UpdateCourseResponse struct {
	Id              uint      `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Description     string    `json:"description"`
	ShortDesc       string    `json:"short_description"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	VideoPreviewURL string    `json:"video_preview_url"`
	Price           float64   `json:"price"`
	DiscountPrice   *float64  `json:"discount_price"`
	InstructorId    uint      `json:"instructor_id"`
	CategoryId      uint      `json:"category_id"`
	CategoryName    string    `json:"category_name"`
	Level           string    `json:"level"`
	DurationHours   int       `json:"duration_hours"`
	TotalLessons    int       `json:"total_lessons"`
	Language        string    `json:"language"`
	Requirements    string    `json:"requirements"`
	WhatYouLearn    string    `json:"what_you_learn"`
	Status          string    `json:"status"`
	IsFeatured      bool      `json:"is_featured"`
	RatingAvg       float32   `json:"rating_avg"`
	RatingCount     int       `json:"rating_count"`
	EnrolledCount   int       `json:"enrolled_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DeleteCourseResponse struct {
	Message  string `json:"message"`
	CourseId uint   `json:"course_id"`
}

type GetCourseStudentsQueryRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Status  string `form:"status" binding:"omitempty,oneof=active completed dropped"`
	Search  string `form:"search" binding:"omitempty,search"`
	OrderBy string `form:"order_by" binding:"omitempty,oneof=enrolled_at completed_at progress_percentage last_accessed_at"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type CourseStudentItem struct {
	UserId             uint       `json:"user_id"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	FullName           string     `json:"full_name"`
	AvatarURL          string     `json:"avatar_url"`
	EnrolledAt         time.Time  `json:"enrolled_at"`
	CompletedAt        *time.Time `json:"completed_at"`
	ProgressPercentage float64    `json:"progress_percentage"`
	LastAccessedAt     *time.Time `json:"last_accessed_at"`
	Status             string     `json:"status"`
}

type GetCourseStudentsResponse struct {
	CourseId    uint                `json:"course_id"`
	CourseTitle string              `json:"course_title"`
	Students    []CourseStudentItem `json:"students"`
	Statistics  StudentStatistics   `json:"statistics"`
	Pagination  PaginationInfo      `json:"pagination"`
}

type StudentStatistics struct {
	TotalStudents     int     `json:"total_students"`
	ActiveStudents    int     `json:"active_students"`
	CompletedStudents int     `json:"completed_students"`
	DroppedStudents   int     `json:"dropped_students"`
	AverageProgress   float64 `json:"average_progress"`
}

type CreateLessonRequest struct {
	Title         string `json:"title" binding:"required,min=3,max=200"`
	Description   string `json:"description" binding:"required,min=10"`
	Content       string `json:"content" binding:"omitempty"`
	VideoURL      string `json:"video_url" binding:"omitempty,url"`
	VideoDuration int    `json:"video_duration" binding:"omitempty,min=0"`
	LessonOrder   int    `json:"lesson_order" binding:"required,min=1"`
	IsPreview     bool   `json:"is_preview" binding:"omitempty"`
	IsPublished   bool   `json:"is_published" binding:"omitempty"`
}

type CreateLessonResponse struct {
	Id            uint   `json:"id"`
	CourseId      uint   `json:"course_id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	VideoURL      string `json:"video_url"`
	VideoDuration int    `json:"video_duration"`
	LessonOrder   int    `json:"lesson_order"`
	IsPreview     bool   `json:"is_preview"`
	IsPublished   bool   `json:"is_published"`
	CreatedAt     string `json:"created_at"`
}

type UpdateLessonRequest struct {
	Title         *string `json:"title" binding:"omitempty,min=3,max=200"`
	Description   *string `json:"description" binding:"omitempty,min=10"`
	Content       *string `json:"content" binding:"omitempty"`
	VideoURL      *string `json:"video_url" binding:"omitempty,url"`
	VideoDuration *int    `json:"video_duration" binding:"omitempty,min=0"`
	LessonOrder   *int    `json:"lesson_order" binding:"omitempty,min=1"`
	IsPreview     *bool   `json:"is_preview" binding:"omitempty"`
	IsPublished   *bool   `json:"is_published" binding:"omitempty"`
}

type UpdateLessonResponse struct {
	Id            uint   `json:"id"`
	CourseId      uint   `json:"course_id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	VideoURL      string `json:"video_url"`
	VideoDuration int    `json:"video_duration"`
	LessonOrder   int    `json:"lesson_order"`
	IsPreview     bool   `json:"is_preview"`
	IsPublished   bool   `json:"is_published"`
	UpdatedAt     string `json:"updated_at"`
}

type DeleteLessonResponse struct {
	Message string `json:"message"`
	Id      uint   `json:"id"`
}

type ReorderLessonsRequest struct {
	Lessons []LessonOrderItem `json:"lessons" binding:"required,min=1,dive"`
}

type LessonOrderItem struct {
	Id          uint `json:"id" binding:"required"`
	LessonOrder int  `json:"lesson_order" binding:"required,min=1"`
}

type ReorderLessonsResponse struct {
	Message      string `json:"message"`
	UpdatedCount int    `json:"updated_count"`
	CourseId     uint   `json:"course_id"`
}
