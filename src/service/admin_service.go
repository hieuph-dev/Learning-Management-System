package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
	"math"
	"strings"
	"time"
)

type adminService struct {
	userRepo   repository.UserRepository
	courseRepo repository.CourseRepository
}

func NewAdminService(userRepo repository.UserRepository, courseRepo repository.CourseRepository) AdminService {
	return &adminService{
		userRepo:   userRepo,
		courseRepo: courseRepo,
	}
}

func (as *adminService) GetUsers(req *dto.GetUsersQueryRequest) (*dto.GetUsersResponse, error) {
	// Set default values

	page := 1
	limit := 10
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
	if req.Role != "" {
		filters["role"] = req.Role
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.Search != "" {
		filters["search"] = utils.NormalizeString(req.Search)
	}

	// Get users with pagination
	users, total, err := as.userRepo.GetUsersWithPagination(offset, limit, filters, orderBy, sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get users", utils.ErrCodeInternal)
	}

	// Convert to DTO
	userItems := make([]dto.AdminUserItem, len(users))
	for i, user := range users {
		userItems[i] = dto.AdminUserItem{
			Id:            user.Id,
			Username:      user.Username,
			Email:         user.Email,
			FullName:      user.FullName,
			Phone:         user.Phone,
			Role:          user.Role,
			Status:        user.Status,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
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

	return &dto.GetUsersResponse{
		Users:      userItems,
		Pagination: pagination,
	}, nil
}

func (as *adminService) GetUserById(userId uint) (*dto.AdminUserDetail, error) {
	// Tìm user theo ID
	user, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// Convert sang DTO
	return &dto.AdminUserDetail{
		Id:            user.Id,
		Username:      user.Username,
		Email:         user.Email,
		FullName:      user.FullName,
		Phone:         user.Phone,
		Bio:           user.Bio,
		AvatarURL:     user.AvatarURL,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

func (as *adminService) UpdateUser(userId uint, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	if existingUser.Status != "active" {
		return nil, utils.NewError("User account is not active", utils.ErrCodeForbidden)
	}

	// 2. Chuẩn bị dữ liệu cập nhật
	updates := make(map[string]interface{})

	if req.FullName != "" {
		updates["full_name"] = strings.TrimSpace(req.FullName)
	}
	if req.Phone != "" {
		updates["phone"] = strings.TrimSpace(req.Phone)
	}
	if req.Bio != "" {
		updates["bio"] = strings.TrimSpace(req.Bio)
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = strings.TrimSpace(req.AvatarURL)
	}
	if req.Role != "" {
		updates["role"] = strings.TrimSpace(req.Role)
	}
	if req.Status != "" {
		updates["status"] = strings.TrimSpace(req.Status)
	}
	// EmailVerified có thể là false nên kiểm tra khác
	updates["email_verified"] = req.EmailVerified

	updates["updated_at"] = time.Now()

	// 3. Cập nhật user
	if err := as.userRepo.UpdateProfile(userId, updates); err != nil {
		return nil, utils.WrapError(err, "Failed to update user", utils.ErrCodeInternal)
	}

	// 4. Lấy thông tin user đã cập nhật
	updatedUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get updated user", utils.ErrCodeInternal)
	}

	return &dto.UpdateUserResponse{
		Id:            updatedUser.Id,
		Username:      updatedUser.Username,
		Email:         updatedUser.Email,
		FullName:      updatedUser.FullName,
		Phone:         updatedUser.Phone,
		Bio:           updatedUser.Bio,
		AvatarURL:     updatedUser.AvatarURL,
		Role:          updatedUser.Role,
		Status:        updatedUser.Status,
		EmailVerified: updatedUser.EmailVerified,
		CreatedAt:     updatedUser.CreatedAt,
		UpdatedAt:     updatedUser.UpdatedAt,
	}, nil
}

func (as *adminService) DeleteUser(userId uint) (*dto.DeleteUserResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// 2. Không cho phép xóa admin khác
	if existingUser.Role == "admin" {
		return nil, utils.NewError("Cannot delete admin account", utils.ErrCodeForbidden)
	}

	// 3. Xóa user khỏi database (soft delete vì model có DeletedAt)
	if err := as.userRepo.DeleteUser(userId); err != nil {
		return nil, utils.WrapError(err, "Failed to delete user", utils.ErrCodeInternal)
	}

	return &dto.DeleteUserResponse{
		Message: "User deleted successfully",
		UserId:  userId,
	}, nil
}

func (as *adminService) ChangeUserStatus(userId uint, req *dto.ChangeUserStatusRequest) (*dto.ChangeUserStatusResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// 2. Không cho phép thay đổi trạng thái admin khác
	if existingUser.Role == "admin" {
		return nil, utils.NewError("Cannot change admin account status", utils.ErrCodeForbidden)
	}

	// 3. Kiểm tra trạng thái hiện tại
	if existingUser.Status == req.Status {
		return nil, utils.NewError(fmt.Sprintf("User is already %s", req.Status), utils.ErrCodeBadRequest)
	}

	// 4. Cập nhật trạng thái
	updates := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}

	if err := as.userRepo.UpdateProfile(userId, updates); err != nil {
		return nil, utils.WrapError(err, "Failed to updated user status", utils.ErrCodeInternal)
	}

	// 5. Tạo message tùy theo trạng thái
	var message string
	switch req.Status {
	case "active":
		message = "User account has been activated"
	case "inactive":
		message = "User account has been deactivated"
	case "banned":
		message = "User account has been banned"
	default:
		message = "User status has been updated"
	}

	if req.Reason != "" {
		message += fmt.Sprintf(". Reason: %s", req.Reason)
	}

	return &dto.ChangeUserStatusResponse{
		Id:       existingUser.Id,
		Username: existingUser.Username,
		Email:    existingUser.Email,
		Status:   req.Status,
		Message:  message,
	}, nil
}

func (as *adminService) GetCourses(req *dto.GetAdminCoursesQueryRequest) (*dto.GetAdminCoursesResponse, error) {
	// Set default values
	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	limit := 20
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := (page - 1) * limit

	// Build filters
	filters := make(map[string]interface{})

	if req.Status != "" {
		filters["status"] = req.Status
	}

	if req.Level != "" {
		filters["level"] = req.Level
	}

	if req.CategoryId > 0 {
		filters["category_id"] = req.CategoryId
	}

	if req.InstructorId > 0 {
		filters["instructor_id"] = req.InstructorId
	}

	if req.IsFeatured != nil {
		filters["is_featured"] = *req.IsFeatured
	}

	if req.Search != "" {
		filters["search"] = req.Search
	}

	// Get courses from repository
	courses, total, err := as.courseRepo.GetCoursesWithPagination(
		offset,
		limit,
		filters,
		req.OrderBy,
		req.SortBy,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get courses", utils.ErrCodeInternal)
	}

	// Convert to DTO
	courseItems := make([]dto.AdminCourseItem, len(courses))
	for i, course := range courses {
		courseItems[i] = dto.AdminCourseItem{
			Id:             course.Id,
			Title:          course.Title,
			Slug:           course.Slug,
			ThumbnailURL:   course.ThumbnailURL,
			Price:          course.Price,
			DiscountPrice:  course.DiscountPrice,
			InstructorId:   course.InstructorId,
			InstructorName: course.Instructor.FullName,
			CategoryId:     course.CategoryId,
			CategoryName:   course.Category.Name,
			Level:          course.Level,
			Status:         course.Status,
			TotalLessons:   course.TotalLessons,
			DurationHours:  course.DurationHours,
			EnrolledCount:  course.EnrolledCount,
			RatingAvg:      course.RatingAvg,
			RatingCount:    course.RatingCount,
			IsFeatured:     course.IsFeatured,
			CreatedAt:      course.CreatedAt,
			UpdatedAt:      course.UpdatedAt,
		}
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.GetAdminCoursesResponse{
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

func (as *adminService) ChangeCourseStatus(courseId uint, req *dto.ChangeCourseStatusRequest) (*dto.ChangeCourseStatusResponse, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := as.courseRepo.FindById(courseId)
	if err != nil {
		return nil, utils.NewError("course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra trạng thái hiện tại
	if course.Status == req.Status {
		return nil, utils.NewError("course already has this status", utils.ErrCodeBadRequest)
	}

	// 3. Validate business rules
	// Không cho phép publish course chưa có lesson
	if req.Status == "published" && course.TotalLessons == 0 {
		return nil, utils.NewError("cannot publish course without lessons", utils.ErrCodeBadRequest)
	}

	// 4. Update status
	if err := as.courseRepo.UpdateCourseStatus(courseId, req.Status); err != nil {
		return nil, utils.WrapError(err, "failed to update course status", utils.ErrCodeInternal)
	}

	// 5. Build message
	message := fmt.Sprintf("Course status changed from '%s' to '%s'", course.Status, req.Status)
	if req.Reason != "" {
		message += fmt.Sprintf(". Reason: %s", req.Reason)
	}

	return &dto.ChangeCourseStatusResponse{
		Id:      course.Id,
		Title:   course.Title,
		Slug:    course.Slug,
		Status:  req.Status,
		Message: message,
	}, nil
}
