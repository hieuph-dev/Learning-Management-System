package service

import (
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
	"math"
)

type courseService struct {
	courseRepo repository.CourseRepository
}

func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{
		courseRepo: courseRepo,
	}
}

func (cs *courseService) GetCourses(req *dto.GetCoursesQueryRequest) (*dto.GetCoursesResponse, error) {
	// Set default values
	page := 1
	limit := 12
	orderBy := "created_at"
	sortBy := "desc"

	if req.Page > 0 {
		page = req.Page
	}

	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}

	if req.SortBy != "" {
		sortBy = req.SortBy
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})

	if req.CategoryId != nil {
		filters["category_id"] = *req.CategoryId
	}

	if req.InstructorId != nil {
		filters["instructor_id"] = *req.InstructorId
	}

	if req.Level != "" {
		filters["level"] = req.Level
	}

	if req.Status != "" {
		filters["status"] = req.Status
	}

	if req.IsFeatured != nil {
		filters["is_featured"] = req.IsFeatured
	}

	if req.Language != "" {
		filters["language"] = req.Language
	}

	if req.MinPrice != nil {
		filters["min_price"] = req.MinPrice
	}

	if req.MaxPrice != nil {
		filters["max_price"] = req.MaxPrice
	}

	if req.Search != "" {
		filters["search"] = req.Search
	}

	// Get courses with pagination
	courses, total, err := cs.courseRepo.GetCoursesWithPagination(offset, limit, filters, orderBy, sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get courses", utils.ErrCodeInternal)
	}

	// Convert to DTO
	courseItems := make([]dto.CourseItem, len(courses))
	for i, course := range courses {
		instructorName := ""
		if course.Instructor.Id != 0 {
			instructorName = course.Instructor.FullName
		}

		categoryName := ""
		if course.Category.Id != 0 {
			categoryName = course.Category.Name
		}

		courseItems[i] = dto.CourseItem{
			Id:             course.Id,
			Title:          course.Title,
			Slug:           course.Slug,
			ShortDesc:      course.ShortDesc,
			ThumbnailURL:   course.ThumbnailURL,
			Price:          course.Price,
			DiscountPrice:  course.DiscountPrice,
			InstructorId:   course.InstructorId,
			InstructorName: instructorName,
			CategoryId:     course.CategoryId,
			CategoryName:   categoryName,
			Level:          course.Level,
			DurationHours:  course.DurationHours,
			TotalLessons:   course.TotalLessons,
			Language:       course.Language,
			Status:         course.Status,
			IsFeatured:     course.IsFeatured,
			RatingAvg:      course.RatingAvg,
			RatingCount:    course.RatingCount,
			EnrolledCount:  course.EnrolledCount,
			CreatedAt:      course.CreatedAt,
		}
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := dto.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	return &dto.GetCoursesResponse{
		Courses:    courseItems,
		Pagination: pagination,
	}, nil

}

func (cs *courseService) SearchCourses(req *dto.SearchCoursesQueryRequest) (*dto.SearchCoursesResponse, error) {
	// Set default values
	page := 1
	limit := 12
	sortBy := "relevance"
	order := "desc"

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 {
		limit = req.Limit
	}
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	if req.Order != "" {
		order = req.Order
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})

	if req.CategoryId != nil {
		filters["category_id"] = *req.CategoryId
	}
	if req.Level != "" {
		filters["level"] = req.Level
	}
	if req.Language != "" {
		filters["language"] = req.Language
	}
	if req.MinPrice != nil {
		filters["min_price"] = *req.MinPrice
	}
	if req.MaxPrice != nil {
		filters["max_price"] = *req.MaxPrice
	}

	// Search courses
	courses, total, err := cs.courseRepo.SearchCourses(req.Q, offset, limit, filters, sortBy, order)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to search courses", utils.ErrCodeInternal)
	}

	// Convert to DTO
	courseItems := make([]dto.CourseItem, len(courses))
	for i, course := range courses {
		instructorName := ""
		if course.Instructor.Id != 0 {
			instructorName = course.Instructor.FullName
		}

		categoryName := ""
		if course.Category.Id != 0 {
			categoryName = course.Category.Name
		}

		courseItems[i] = dto.CourseItem{
			Id:             course.Id,
			Title:          course.Title,
			Slug:           course.Slug,
			ShortDesc:      course.ShortDesc,
			ThumbnailURL:   course.ThumbnailURL,
			Price:          course.Price,
			DiscountPrice:  course.DiscountPrice,
			InstructorId:   course.InstructorId,
			InstructorName: instructorName,
			CategoryId:     course.CategoryId,
			CategoryName:   categoryName,
			Level:          course.Level,
			DurationHours:  course.DurationHours,
			TotalLessons:   course.TotalLessons,
			Language:       course.Language,
			Status:         course.Status,
			IsFeatured:     course.IsFeatured,
			RatingAvg:      course.RatingAvg,
			RatingCount:    course.RatingCount,
			EnrolledCount:  course.EnrolledCount,
			CreatedAt:      course.CreatedAt,
		}
	}

	// Get search filters
	searchFilters, err := cs.courseRepo.GetSearchFilters(req.Q)
	if err != nil {
		// Log error but don't fail the request
		searchFilters = &dto.SearchFilters{}
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := dto.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return &dto.SearchCoursesResponse{
		Query:      req.Q,
		Courses:    courseItems,
		Pagination: pagination,
		Filters:    *searchFilters,
	}, nil
}

func (cs *courseService) GetFeaturedCourses(req *dto.GetFeaturedCoursesQueryRequest) (*dto.GetFeaturedCoursesResponse, error) {
	// Set default limit
	limit := 8
	if req.Limit > 0 {
		limit = req.Limit
	}

	// Prepare filters
	filters := make(map[string]interface{})

	if req.CategoryId != nil {
		filters["category_id"] = *req.CategoryId
	}
	if req.Level != "" {
		filters["level"] = req.Level
	}
	if req.Language != "" {
		filters["language"] = req.Language
	}

	// Get featured courses
	courses, total, err := cs.courseRepo.GetFeaturedCourses(limit, filters)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get featured courses", utils.ErrCodeInternal)
	}

	// Convert to DTO
	courseItems := make([]dto.CourseItem, len(courses))
	for i, course := range courses {
		instructorName := ""
		if course.Instructor.Id != 0 {
			instructorName = course.Instructor.FullName
		}

		categoryName := ""
		if course.Category.Id != 0 {
			categoryName = course.Category.Name
		}

		courseItems[i] = dto.CourseItem{
			Id:             course.Id,
			Title:          course.Title,
			Slug:           course.Slug,
			ShortDesc:      course.ShortDesc,
			ThumbnailURL:   course.ThumbnailURL,
			Price:          course.Price,
			DiscountPrice:  course.DiscountPrice,
			InstructorId:   course.InstructorId,
			InstructorName: instructorName,
			CategoryId:     course.CategoryId,
			CategoryName:   categoryName,
			Level:          course.Level,
			DurationHours:  course.DurationHours,
			TotalLessons:   course.TotalLessons,
			Language:       course.Language,
			Status:         course.Status,
			IsFeatured:     course.IsFeatured,
			RatingAvg:      course.RatingAvg,
			RatingCount:    course.RatingCount,
			EnrolledCount:  course.EnrolledCount,
			CreatedAt:      course.CreatedAt,
		}
	}

	return &dto.GetFeaturedCoursesResponse{
		Courses: courseItems,
		Total:   total,
	}, nil

}

func (cs *courseService) GetCourseBySlug(slug string) (*dto.CourseDetail, error) {
	// Tìm course theo slug
	course, err := cs.courseRepo.FindBySlug(slug)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// Kiểm tra trạng thái course - chỉ hiển thị course đã published
	if course.Status != "published" {
		return nil, utils.NewError("Course not available", utils.ErrCodeNotFound)
	}

	// Get instructor info
	instructorName := ""
	instructorBio := ""
	if course.Instructor.Id != 0 {
		instructorName = course.Instructor.FullName
		instructorBio = course.Instructor.Bio
	}

	// Get category info
	categoryName := ""
	if course.Category.Id != 0 {
		categoryName = course.Category.Name
	}

	return &dto.CourseDetail{
		Id:              course.Id,
		Title:           course.Title,
		Slug:            course.Slug,
		Description:     course.Description,
		ShortDesc:       course.ShortDesc,
		ThumbnailURL:    course.ThumbnailURL,
		VideoPreviewURL: course.VideoPreviewURL,
		Price:           course.Price,
		DiscountPrice:   course.DiscountPrice,
		InstructorId:    course.InstructorId,
		InstructorName:  instructorName,
		InstructorBio:   instructorBio,
		CategoryId:      course.CategoryId,
		CategoryName:    categoryName,
		Level:           course.Level,
		DurationHours:   course.DurationHours,
		TotalLessons:    course.TotalLessons,
		Language:        course.Language,
		Requirements:    course.Requirements,
		WhatYouLearn:    course.WhatYouLearn,
		Status:          course.Status,
		IsFeatured:      course.IsFeatured,
		RatingAvg:       course.RatingAvg,
		RatingCount:     course.RatingCount,
		EnrolledCount:   course.EnrolledCount,
		CreatedAt:       course.CreatedAt,
		UpdatedAt:       course.UpdatedAt,
	}, nil
}
