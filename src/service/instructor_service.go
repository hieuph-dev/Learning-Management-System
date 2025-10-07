package service

import (
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"math"
	"strconv"
)

type instructorService struct {
	instructorRepo repository.InstructorRepository
	categoryRepo   repository.CategoryRepository
}

func NewInstructorService(instructorRepo repository.InstructorRepository, categoryRepo repository.CategoryRepository) InstructorService {
	return &instructorService{
		instructorRepo: instructorRepo,
		categoryRepo:   categoryRepo,
	}
}

func (is *instructorService) GetInstructorCourses(instructorId uint, req *dto.GetInstructorCoursesQueryRequest) (*dto.GetInstructorCoursesResponse, error) {
	// Set default values
	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := (page - 1) * limit

	// Build filters
	filters := make(map[string]interface{})

	if req.Status != "" {
		filters["status"] = req.Status
	}

	if req.Search != "" {
		filters["search"] = req.Search
	}

	// Get courses from repository
	courses, total, err := is.instructorRepo.GetInstructorCourses(
		instructorId,
		offset,
		limit,
		filters,
		req.OrderBy,
		req.SortBy,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get instructor courses", utils.ErrCodeInternal)
	}

	// Convert to DTO
	courseItems := make([]dto.InstructorCourseItem, len(courses))
	for i, course := range courses {
		courseItems[i] = dto.InstructorCourseItem{
			Id:            course.Id,
			Title:         course.Title,
			Slug:          course.Slug,
			ThumbnailURL:  course.ThumbnailURL,
			Price:         course.Price,
			DiscountPrice: course.DiscountPrice,
			CategoryId:    course.CategoryId,
			CategoryName:  course.Category.Name,
			Level:         course.Level,
			Status:        course.Status,
			TotalLessons:  course.TotalLessons,
			DurationHours: course.DurationHours,
			EnrolledCount: course.EnrolledCount,
			RatingAvg:     course.RatingAvg,
			RatingCount:   course.RatingCount,
			IsFeatured:    course.IsFeatured,
			CreatedAt:     course.CreatedAt,
			UpdatedAt:     course.UpdatedAt,
		}
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.GetInstructorCoursesResponse{
		Courses: courseItems,
		Pagination: dto.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

func (is *instructorService) CreateCourse(instructorId uint, req *dto.CreateCourseRequest) (*dto.CreateCourseResponse, error) {
	// 1. Validate category exists và active
	category, err := is.categoryRepo.FindById(req.CategoryId)
	if err != nil {
		return nil, utils.NewError("category not found", utils.ErrCodeNotFound)
	}

	if !category.IsActive {
		return nil, utils.NewError("category is not active", utils.ErrCodeBadRequest)
	}

	// 2. Validate discount price
	if req.DiscountPrice != nil && *req.DiscountPrice >= req.Price {
		return nil, utils.NewError("discount price must be less than regular price", utils.ErrCodeBadRequest)
	}

	// 3. Generate unique slug
	baseSlug := utils.GenerateSlug(req.Title)
	uniqueSlug := utils.GenerateUniqueSlug(baseSlug, func(slug string) bool {
		_, exists := is.instructorRepo.FindCourseBySlug(slug)
		return exists
	})

	// 4. Create course model
	course := &models.Course{
		Title:         req.Title,
		Slug:          uniqueSlug,
		Description:   req.Description,
		ShortDesc:     req.ShortDesc,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		InstructorId:  instructorId,
		CategoryId:    req.CategoryId,
		Level:         req.Level,
		Language:      req.Language,
		Requirements:  req.Requirements,
		WhatYouLearn:  req.WhatYouLearn,
		DurationHours: req.DurationHours,
		Status:        "draft", // Mặc định là draft
		IsFeatured:    false,
		TotalLessons:  0,
		RatingAvg:     0,
		RatingCount:   0,
		EnrolledCount: 0,
	}

	// 5. Save to database
	if err := is.instructorRepo.CreateCourse(course); err != nil {
		return nil, utils.WrapError(err, "failed to create course", utils.ErrCodeInternal)
	}

	// 6. Return response
	return &dto.CreateCourseResponse{
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
		CategoryId:      course.CategoryId,
		CategoryName:    category.Name,
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

func (is *instructorService) UpdateCourse(instructorId, courseId uint, req *dto.UpdateCourseRequest) (*dto.UpdateCourseResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor này không
	course, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("course not found or you don't have permission to update this course", utils.ErrCodeNotFound)
	}

	// 2. Validate category nếu được cập nhật
	var category *models.Category
	if req.CategoryId != 0 && req.CategoryId != course.CategoryId {
		category, err = is.categoryRepo.FindById(req.CategoryId)
		if err != nil {
			return nil, utils.NewError("category not found", utils.ErrCodeNotFound)
		}
		if !category.IsActive {
			return nil, utils.NewError("category is not active", utils.ErrCodeBadRequest)
		}
	} else {
		// Giữ nguyên category hiện tại
		category, _ = is.categoryRepo.FindById(course.CategoryId)
	}

	// 3. Validate discount price nếu có
	finalPrice := course.Price
	if req.Price > 0 {
		finalPrice = req.Price
	}
	if req.DiscountPrice != nil && *req.DiscountPrice >= finalPrice {
		return nil, utils.NewError("discount price must be less than regular price", utils.ErrCodeBadRequest)
	}

	// 4. Build updates map
	updates := make(map[string]interface{})

	if req.Title != "" && req.Title != course.Title {
		updates["title"] = req.Title
		// Generate new slug nếu title thay đổi
		baseSlug := utils.GenerateSlug(req.Title)
		uniqueSlug := utils.GenerateUniqueSlug(baseSlug, func(slug string) bool {
			if slug == course.Slug {
				return false // Cho phép giữ nguyên slug hiện tại
			}
			_, exists := is.instructorRepo.FindCourseBySlug(slug)
			return exists
		})
		updates["slug"] = uniqueSlug
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}

	if req.ShortDesc != "" {
		updates["short_desc"] = req.ShortDesc
	}

	if req.CategoryId != 0 {
		updates["category_id"] = req.CategoryId
	}

	if req.Level != "" {
		updates["level"] = req.Level
	}

	if req.Language != "" {
		updates["language"] = req.Language
	}

	if req.Price > 0 {
		updates["price"] = req.Price
	}

	if req.DiscountPrice != nil {
		updates["discount_price"] = req.DiscountPrice
	}

	if req.Requirements != "" {
		updates["requirements"] = req.Requirements
	}

	if req.WhatYouLearn != "" {
		updates["what_you_learn"] = req.WhatYouLearn
	}

	if req.DurationHours >= 0 {
		updates["duration_hours"] = req.DurationHours
	}

	if req.Status != "" {
		// Không cho phép chuyển từ published về draft nếu đã có học viên
		if course.Status == "published" && req.Status == "draft" {
			enrollmentCount, _ := is.instructorRepo.CountEnrollmentsByCourse(courseId)
			if enrollmentCount > 0 {
				return nil, utils.NewError("cannot change status to draft when course has enrollments", utils.ErrCodeBadRequest)
			}
		}
		updates["status"] = req.Status
	}

	if req.IsFeatured != nil {
		updates["is_featured"] = *req.IsFeatured
	}

	// 5. Kiểm tra có gì cần update không
	if len(updates) == 0 {
		return nil, utils.NewError("no fields to update", utils.ErrCodeBadRequest)
	}

	// 6. Update course
	if err := is.instructorRepo.UpdateCourse(courseId, updates); err != nil {
		return nil, utils.WrapError(err, "failed to update course", utils.ErrCodeInternal)
	}

	// 7. Lấy lại course đã update
	updatedCourse, err := is.instructorRepo.FindCourseById(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "failed to fetch updated course", utils.ErrCodeInternal)
	}

	// 8. Return response
	return &dto.UpdateCourseResponse{
		Id:              updatedCourse.Id,
		Title:           updatedCourse.Title,
		Slug:            updatedCourse.Slug,
		Description:     updatedCourse.Description,
		ShortDesc:       updatedCourse.ShortDesc,
		ThumbnailURL:    updatedCourse.ThumbnailURL,
		VideoPreviewURL: updatedCourse.VideoPreviewURL,
		Price:           updatedCourse.Price,
		DiscountPrice:   updatedCourse.DiscountPrice,
		InstructorId:    updatedCourse.InstructorId,
		CategoryId:      updatedCourse.CategoryId,
		CategoryName:    category.Name,
		Level:           updatedCourse.Level,
		DurationHours:   updatedCourse.DurationHours,
		TotalLessons:    updatedCourse.TotalLessons,
		Language:        updatedCourse.Language,
		Requirements:    updatedCourse.Requirements,
		WhatYouLearn:    updatedCourse.WhatYouLearn,
		Status:          updatedCourse.Status,
		IsFeatured:      updatedCourse.IsFeatured,
		RatingAvg:       updatedCourse.RatingAvg,
		RatingCount:     updatedCourse.RatingCount,
		EnrolledCount:   updatedCourse.EnrolledCount,
		CreatedAt:       updatedCourse.CreatedAt,
		UpdatedAt:       updatedCourse.UpdatedAt,
	}, nil
}

func (is *instructorService) DeleteCourse(instructorId, courseId uint) (*dto.DeleteCourseResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor này không
	course, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("course not found or you don't have permission to delete this course", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra xem course có học viên đang học không
	enrollmentCount, err := is.instructorRepo.CountEnrollmentsByCourse(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "failed to check enrollments", utils.ErrCodeInternal)
	}

	if enrollmentCount > 0 {
		return nil, utils.NewError("cannot delete course with active enrollments", utils.ErrCodeBadRequest)
	}

	// 3. Chỉ cho phép xóa course ở trạng thái draft hoặc archived
	if course.Status == "published" {
		return nil, utils.NewError("cannot delete published course. Please archive it first", utils.ErrCodeBadRequest)
	}

	// 4. Xóa course (soft delete)
	if err := is.instructorRepo.DeleteCourse(courseId); err != nil {
		return nil, utils.WrapError(err, "failed to delete course", utils.ErrCodeInternal)
	}

	return &dto.DeleteCourseResponse{
		Message:  "Course deleted successfully",
		CourseId: courseId,
	}, nil
}

func (is *instructorService) GetCourseStudents(instructorId, courseId uint, req *dto.GetCourseStudentsQueryRequest) (*dto.GetCourseStudentsResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor này không
	course, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("course not found or you don't have permission to view students", utils.ErrCodeNotFound)
	}

	// 2. Set default values
	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	limit := 20
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := (page - 1) * limit

	// 3. Build filters
	filters := make(map[string]interface{})

	if req.Status != "" {
		filters["status"] = req.Status
	}

	if req.Search != "" {
		filters["search"] = req.Search
	}

	// 4. Get enrollments from repository
	enrollments, total, err := is.instructorRepo.GetCourseStudents(
		courseId,
		offset,
		limit,
		filters,
		req.OrderBy,
		req.SortBy,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get course students", utils.ErrCodeInternal)
	}

	// 5. Get statistics
	statistics, err := is.instructorRepo.GetStudentStatistics(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get student statistics", utils.ErrCodeInternal)
	}

	// 6. Convert to DTO
	studentItems := make([]dto.CourseStudentItem, len(enrollments))
	for i, enrollment := range enrollments {
		studentItems[i] = dto.CourseStudentItem{
			UserId:             enrollment.UserId,
			Username:           enrollment.User.Username,
			Email:              enrollment.User.Email,
			FullName:           enrollment.User.FullName,
			AvatarURL:          enrollment.User.AvatarURL,
			EnrolledAt:         enrollment.EnrolledAt,
			CompletedAt:        enrollment.CompletedAt,
			ProgressPercentage: enrollment.ProgressPercentage,
			LastAccessedAt:     enrollment.LastAccessedAt,
			Status:             enrollment.Status,
		}
	}

	// 7. Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.GetCourseStudentsResponse{
		CourseId:    course.Id,
		CourseTitle: course.Title,
		Students:    studentItems,
		Statistics:  *statistics,
		Pagination: dto.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

func (is *instructorService) CreateLesson(instructorId, courseId uint, req *dto.CreateLessonRequest) (*dto.CreateLessonResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor không
	_, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("Course not found or you don't have permission", utils.ErrCodeNotFound)
	}

	// 2. Generate slug từ title
	baseSlug := utils.GenerateSlug(req.Title)

	// 3. Kiểm tra slug đã tồn tại chưa
	slug := baseSlug
	counter := 1
	for {
		_, exists := is.instructorRepo.FindLessonBySlug(slug, courseId)
		if !exists {
			break
		}
		slug = baseSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	// 4. Kiểm tra lesson_order đã tồn tại chưa
	orderExists, err := is.instructorRepo.CheckLessonOrderExists(courseId, req.LessonOrder)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to check lesson order", utils.ErrCodeInternal)
	}
	if orderExists {
		return nil, utils.NewError("Lesson order already exists in this course", utils.ErrCodeConflict)
	}

	// 5. Tạo lesson mới
	lesson := &models.Lesson{
		CourseId:      courseId,
		Title:         req.Title,
		Slug:          slug,
		Description:   req.Description,
		Content:       req.Content,
		VideoURL:      req.VideoURL,
		VideoDuration: req.VideoDuration,
		LessonOrder:   req.LessonOrder,
		IsPreview:     req.IsPreview,
		IsPublished:   req.IsPublished,
	}

	// 6. Lưu vào database
	if err := is.instructorRepo.CreateLesson(lesson); err != nil {
		return nil, utils.WrapError(err, "Failed to create lesson", utils.ErrCodeInternal)
	}

	// 7. Trả về response
	return &dto.CreateLessonResponse{
		Id:            lesson.Id,
		CourseId:      lesson.CourseId,
		Title:         lesson.Title,
		Slug:          lesson.Slug,
		Description:   lesson.Description,
		Content:       lesson.Content,
		VideoURL:      lesson.VideoURL,
		VideoDuration: lesson.VideoDuration,
		LessonOrder:   lesson.LessonOrder,
		IsPreview:     lesson.IsPreview,
		IsPublished:   lesson.IsPublished,
		CreatedAt:     lesson.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (is *instructorService) UpdateLesson(instructorId, courseId, lessonId uint, req *dto.UpdateLessonRequest) (*dto.UpdateLessonResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor không
	_, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("Course not found or you don't have permission", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra lesson có tồn tại và thuộc về course không
	_, err = is.instructorRepo.FindLessonByIdAndCourse(lessonId, courseId)
	if err != nil {
		return nil, utils.NewError("Lesson not found", utils.ErrCodeNotFound)
	}

	// 3. Chuẩn bị updates map
	updates := make(map[string]interface{})

	if req.Title != nil {
		updates["title"] = *req.Title
		// Generate slug mới nếu title thay đổi
		baseSlug := utils.GenerateSlug(*req.Title)
		slug := baseSlug
		counter := 1
		for {
			existingLesson, exists := is.instructorRepo.FindLessonBySlug(slug, courseId)
			if !exists || existingLesson.Id == lessonId {
				break
			}
			slug = baseSlug + "-" + strconv.Itoa(counter)
			counter++
		}
		updates["slug"] = slug
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Content != nil {
		updates["content"] = *req.Content
	}

	if req.VideoURL != nil {
		updates["video_url"] = *req.VideoURL
	}

	if req.VideoDuration != nil {
		updates["video_duration"] = *req.VideoDuration
	}

	if req.LessonOrder != nil {
		// Kiểm tra lesson_order mới có bị trùng không (trừ chính lesson này)
		orderExists, err := is.instructorRepo.CheckLessonOrderExistsExcept(courseId, *req.LessonOrder, lessonId)
		if err != nil {
			return nil, utils.WrapError(err, "Failed to check lesson order", utils.ErrCodeInternal)
		}
		if orderExists {
			return nil, utils.NewError("Lesson order already exists in this course", utils.ErrCodeConflict)
		}
		updates["lesson_order"] = *req.LessonOrder
	}

	if req.IsPreview != nil {
		updates["is_preview"] = *req.IsPreview
	}

	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}

	// 4. Nếu không có gì để update
	if len(updates) == 0 {
		return nil, utils.NewError("No fields to update", utils.ErrCodeBadRequest)
	}

	// 5. Update lesson
	if err := is.instructorRepo.UpdateLesson(lessonId, updates); err != nil {
		return nil, utils.WrapError(err, "Failed to update lesson", utils.ErrCodeInternal)
	}

	// 6. Lấy lại lesson đã update
	updatedLesson, err := is.instructorRepo.FindLessonByIdAndCourse(lessonId, courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get updated lesson", utils.ErrCodeInternal)
	}

	// 7. Trả về response
	return &dto.UpdateLessonResponse{
		Id:            updatedLesson.Id,
		CourseId:      updatedLesson.CourseId,
		Title:         updatedLesson.Title,
		Slug:          updatedLesson.Slug,
		Description:   updatedLesson.Description,
		Content:       updatedLesson.Content,
		VideoURL:      updatedLesson.VideoURL,
		VideoDuration: updatedLesson.VideoDuration,
		LessonOrder:   updatedLesson.LessonOrder,
		IsPreview:     updatedLesson.IsPreview,
		IsPublished:   updatedLesson.IsPublished,
		UpdatedAt:     updatedLesson.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (is *instructorService) DeleteLesson(instructorId, courseId, lessonId uint) (*dto.DeleteLessonResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor không
	_, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("Course not found or you don't have permission", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra lesson có tồn tại và thuộc về course không
	_, err = is.instructorRepo.FindLessonByIdAndCourse(lessonId, courseId)
	if err != nil {
		return nil, utils.NewError("Lesson not found", utils.ErrCodeNotFound)
	}

	// 3. Delete lesson (soft delete)
	if err := is.instructorRepo.DeleteCourse(lessonId); err != nil {
		return nil, utils.WrapError(err, "Failed to delete lesson", utils.ErrCodeInternal)
	}

	// 4. Trả về response
	return &dto.DeleteLessonResponse{
		Message: "Lesson deleted successfully",
		Id:      lessonId,
	}, nil
}

func (is *instructorService) ReorderLessons(instructorId, lessonId uint, req *dto.ReorderLessonsRequest) (*dto.ReorderLessonsResponse, error) {
	// 1. Lấy thông tin lesson để biết course_id
	// (Sử dụng bất kỳ lesson nào trong danh sách để lấy course_id)
	if len(req.Lessons) == 0 {
		return nil, utils.NewError("Lessons list cannot be empty", utils.ErrCodeBadRequest)
	}

	// Lấy lesson đầu tiên để xác định course
	// firstLessonId := req.Lessons[0].Id
	// var firstLesson models.Lesson

	// Tìm lesson bất kỳ để lấy course_id
	lessonIds := make([]uint, len(req.Lessons))
	for i, item := range req.Lessons {
		lessonIds[i] = item.Id
	}

	lessons, err := is.instructorRepo.FindLessonsByIds(lessonIds)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to find lessons", utils.ErrCodeInternal)
	}

	if len(lessons) == 0 {
		return nil, utils.NewError("No lessons found", utils.ErrCodeNotFound)
	}

	// Lấy course_id từ lesson đầu tiên
	courseId := lessons[0].CourseId
	// firstLesson = lessons[0]

	// 2. Kiểm tra course có tồn tại và thuộc về instructor không
	_, err = is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("Course not found or you don't have permission", utils.ErrCodeNotFound)
	}

	// 3. Kiểm tra tất cả lessons có thuộc về cùng một course không
	lessonCourseMap := make(map[uint]uint) // lessonId -> courseId
	for _, lesson := range lessons {
		lessonCourseMap[lesson.Id] = lesson.CourseId
		if lesson.CourseId != courseId {
			return nil, utils.NewError("All lessons must belong to the same course", utils.ErrCodeBadRequest)
		}
	}

	// 4. Kiểm tra số lượng lessons trong request có khớp với số lessons tìm được không
	if len(lessons) != len(req.Lessons) {
		return nil, utils.NewError("Some lessons not found or already deleted", utils.ErrCodeNotFound)
	}

	// 5. Kiểm tra không có lesson_order trùng nhau
	orderMap := make(map[int]bool)
	for _, item := range req.Lessons {
		if orderMap[item.LessonOrder] {
			return nil, utils.NewError("Duplicate lesson orders found", utils.ErrCodeBadRequest)
		}
		orderMap[item.LessonOrder] = true
	}

	// 6. Bắt đầu transaction để update
	tx := is.instructorRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 7. Update lesson orders
	updateCount := 0
	for _, item := range req.Lessons {
		err := tx.Model(&models.Lesson{}).
			Where("id = ?", item.Id).
			Update("lesson_order", item.LessonOrder).Error
		if err != nil {
			tx.Rollback()
			return nil, utils.WrapError(err, "Failed to update lesson order", utils.ErrCodeInternal)
		}
		updateCount++
	}

	// 8. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, utils.WrapError(err, "Failed to commit transaction", utils.ErrCodeInternal)
	}

	// 9. Trả về response
	return &dto.ReorderLessonsResponse{
		Message:      "Lessons reordered successfully",
		UpdatedCount: updateCount,
		CourseId:     courseId,
	}, nil

}
