package dto

import "time"

type CategoryItem struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ParentId    *uint     `json:"parent_id"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetCategoriesResponse struct {
	Categories []CategoryItem `json:"categories"`
	Total      int            `json:"total"`
}

type GetCategoriesQueryRequest struct {
	ParentId *uint  `form:"parent_id" binding:"omitempty"`
	IsActive *bool  `form:"is_active" binding:"omitempty"`
	Search   string `form:"search" binding:"omitempty,search"`
}

type CategoryDetail struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ParentId    *uint     `json:"parent_id"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Slug        string `json:"slug" binding:"required,slug,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	ParentId    *uint  `json:"parent_id" binding:"omitempty"`
	SortOrder   int    `json:"sort_order" binding:"omitempty,min=0"`
	IsActive    bool   `json:"is_active"`
}

type CreateCategoryResponse struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ParentId    *uint     `json:"parent_id"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Slug        string `json:"slug" binding:"omitempty,slug,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	ParentId    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateCategoryResponse struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ParentId    *uint     `json:"parent_id"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeleteCategoryResponse struct {
	Message    string `json:"message"`
	CategoryId uint   `json:"category_id"`
}
