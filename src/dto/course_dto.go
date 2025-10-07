package dto

import "time"

type CourseItem struct {
	Id             uint      `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	ShortDesc      string    `json:"short_description"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	Price          float64   `json:"price"`
	DiscountPrice  *float64  `json:"discount_price"`
	InstructorId   uint      `json:"instructor_id"`
	InstructorName string    `json:"instructor_name"`
	CategoryId     uint      `json:"category_id"`
	CategoryName   string    `json:"category_name"`
	Level          string    `json:"level"`
	DurationHours  int       `json:"duration_hours"`
	TotalLessons   int       `json:"total_lessons"`
	Language       string    `json:"language"`
	Status         string    `json:"status"`
	IsFeatured     bool      `json:"is_featured"`
	RatingAvg      float32   `json:"rating_avg"`
	RatingCount    int       `json:"rating_count"`
	EnrolledCount  int       `json:"enrolled_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetCoursesQueryRequest struct {
	Page         int      `form:"page" binding:"omitempty,min=1"`
	Limit        int      `form:"limit" binding:"omitempty,min=1,max=100"`
	CategoryId   *uint    `form:"category_id" binding:"omitempty"`
	InstructorId *uint    `form:"instructor_id" binding:"omitempty"`
	Level        string   `form:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Status       string   `form:"status" binding:"omitempty,oneof=draft published archived"`
	IsFeatured   *bool    `form:"is_featured" binding:"omitempty"`
	Language     string   `form:"language" binding:"omitempty,oneof=vi en"`
	MinPrice     *float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice     *float64 `form:"max_price" binding:"omitempty,min=0"`
	Search       string   `form:"search" binding:"omitempty,search"`
	OrderBy      string   `form:"order_by" binding:"omitempty,oneof=created_at updated_at price rating_avg enrolled_count"`
	SortBy       string   `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type GetCoursesResponse struct {
	Courses    []CourseItem   `json:"courses"`
	Pagination PaginationInfo `json:"pagination"`
}

type SearchCoursesQueryRequest struct {
	Q          string   `form:"q" binding:"required,min=2"`
	Page       int      `form:"page" binding:"omitempty,min=1"`
	Limit      int      `form:"limit" binding:"omitempty,min=1,max=50"`
	CategoryId *uint    `form:"category_id" binding:"omitempty"`
	Level      string   `form:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	MinPrice   *float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice   *float64 `form:"max_price" binding:"omitempty,min=0"`
	Language   string   `form:"language" binding:"omitempty,oneof=vi en"`
	SortBy     string   `form:"sort_by" binding:"omitempty,oneof=relevance price rating_avg enrolled_count created_at"`
	Order      string   `form:"order" binding:"omitempty,oneof=asc desc"`
}

type SearchCoursesResponse struct {
	Query      string         `json:"query"`
	Courses    []CourseItem   `json:"courses"`
	Pagination PaginationInfo `json:"pagination"`
	Filters    SearchFilters  `json:"filters"`
}

type SearchFilters struct {
	Categories  []FilterOption `json:"categories"`
	Levels      []FilterOption `json:"levels"`
	PriceRanges []FilterOption `json:"price_ranges"`
	Languages   []FilterOption `json:"languages"`
}

type FilterOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type GetFeaturedCoursesQueryRequest struct {
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=20"`
	CategoryId *uint  `form:"category_id" binding:"omitempty"`
	Level      string `form:"level" binding:"omitempty,oneof=beginner intermediate advanced"`
	Language   string `form:"language" binding:"omitempty,oneof=vi en"`
}

type GetFeaturedCoursesResponse struct {
	Courses []CourseItem `json:"courses"`
	Total   int          `json:"total"`
}

type CourseDetail struct {
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
	InstructorName  string    `json:"instructor_name"`
	InstructorBio   string    `json:"instructor_bio"`
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

type ReviewItem struct {
	Id          uint      `json:"id"`
	UserId      uint      `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserAvatar  string    `json:"user_avatar"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetCourseReviewsQueryRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=50"`
	Rating    *int   `form:"rating" binding:"omitempty,min=1,max=5"`
	Published *bool  `form:"published" binding:"omitempty"`
	OrderBy   string `form:"order_by" binding:"omitempty,oneof=created_at rating"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type GetCourseReviewsResponse struct {
	Reviews    []ReviewItem   `json:"reviews"`
	Pagination PaginationInfo `json:"pagination"`
	Stats      ReviewStats    `json:"stats"`
}

type ReviewStats struct {
	TotalReviews       int         `json:"total_reviews"`
	AverageRating      float64     `json:"average_rating"`
	RatingDistribution map[int]int `json:"rating_distribution"`
}
