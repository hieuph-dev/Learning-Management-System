package service

import (
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"math"
)

type reviewService struct {
	reviewRepo     repository.ReviewRepository
	courseRepo     repository.CourseRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, courseRepo repository.CourseRepository, enrollmentRepo repository.EnrollmentRepository) ReviewService {
	return &reviewService{
		reviewRepo:     reviewRepo,
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (rs *reviewService) GetCourseReviews(courseId uint, req *dto.GetCourseReviewsQueryRequest) (*dto.GetCourseReviewsResponse, error) {
	// Verify course exists
	_, err := rs.courseRepo.FindBySlug("")
	if err != nil {
		// For now, skip course verification or add FindByID to CourseRepository
	}

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

	if req.Rating != nil {
		filters["rating"] = *req.Rating
	}
	if req.Published != nil {
		filters["is_published"] = *req.Published
	} else {
		// Default: only show published reviews
		filters["is_published"] = true
	}

	// Get reviews with pagination
	reviews, total, err := rs.reviewRepo.GetCourseReviews(courseId, offset, limit, filters, orderBy, sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get course reviews", utils.ErrCodeInternal)
	}

	// Convert to DTO
	reviewItems := make([]dto.ReviewItem, len(reviews))
	for i, review := range reviews {
		userName := "Annymous"
		userAvatar := ""

		if review.User.Id != 0 {
			userName = review.User.FullName
			userAvatar = review.User.AvatarURL
		}

		reviewItems[i] = dto.ReviewItem{
			Id:          review.Id,
			UserId:      review.UserId,
			UserName:    userName,
			UserAvatar:  userAvatar,
			Rating:      review.Rating,
			Comment:     review.Comment,
			IsPublished: review.IsPublished,
			CreatedAt:   review.CreatedAt,
		}
	}

	// Get review stats
	stats, err := rs.reviewRepo.GetCourseReviewStats(courseId)
	if err != nil {
		// Log error but don't fail the request
		stats = &dto.ReviewStats{
			RatingDistribution: make(map[int]int),
		}
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

	return &dto.GetCourseReviewsResponse{
		Reviews:    reviewItems,
		Pagination: pagination,
		Stats:      *stats,
	}, nil
}

func (rs *reviewService) CreateReview(userId, courseId uint, req *dto.CreateReviewRequest) (*dto.CreateReviewResponse, error) {
	// Check if course exists
	_, err := rs.courseRepo.FindById(courseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// Check if user enrolled the course
	isEnrolled, err := rs.enrollmentRepo.CheckUserEnrollment(userId, courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to check enrollment", utils.ErrCodeInternal)
	}
	if !isEnrolled {
		return nil, utils.NewError("You must be enrolled in this course to write a review", utils.ErrCodeForbidden)
	}

	// Check if user already reviewed this course
	existtingReview, err := rs.reviewRepo.FindByUserAndCourse(userId, courseId)
	if err == nil && existtingReview != nil {
		return nil, utils.NewError("You have already reviewed this course", utils.ErrCodeConflict)
	}

	// Create review
	review := &models.Review{
		UserId:      userId,
		CourseId:    courseId,
		Rating:      req.Rating,
		Comment:     req.Comment,
		IsPublished: true,
	}

	if err := rs.reviewRepo.Create(review); err != nil {
		return nil, utils.WrapError(err, "Failed to create review", utils.ErrCodeInternal)
	}

	// Update course rating stats
	if err := rs.reviewRepo.UpdateCourseRatingStats(courseId); err != nil {
		// Log error but don't fail the request
	}

	return &dto.CreateReviewResponse{
		Id:       review.Id,
		CourseId: courseId,
		Rating:   review.Rating,
		Comment:  review.Comment,
		Message:  "Review created successfully",
	}, nil
}

func (rs *reviewService) UpdateReview(userId, reviewId uint, req *dto.UpdateReviewRequest) (*dto.UpdateReviewResponse, error) {
	// Find review
	review, err := rs.reviewRepo.FindById(reviewId)
	if err != nil {
		return nil, utils.NewError("Review not found", utils.ErrCodeNotFound)
	}

	// Check ownership
	if review.UserId != userId {
		return nil, utils.NewError("You can only update your own review", utils.ErrCodeForbidden)
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if req.Rating != nil {
		updates["rating"] = *req.Rating
		review.Rating = *req.Rating
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
		review.Comment = *req.Comment
	}

	// Update review
	if err := rs.reviewRepo.Update(reviewId, updates); err != nil {
		return nil, utils.WrapError(err, "Failed to update review", utils.ErrCodeInternal)
	}

	// Update course rating stats
	if req.Rating != nil {
		if err := rs.reviewRepo.UpdateCourseRatingStats(review.CourseId); err != nil {
			// Log error but don't fail the request
		}
	}

	return &dto.UpdateReviewResponse{
		Id:      review.Id,
		Rating:  review.Rating,
		Comment: review.Comment,
		Message: "Review updated successfully",
	}, nil
}

func (rs *reviewService) DeleteReview(userId, reviewId uint) (*dto.DeleteReviewResponse, error) {
	// Find review
	review, err := rs.reviewRepo.FindById(reviewId)
	if err != nil {
		return nil, utils.NewError("Review not found", utils.ErrCodeNotFound)
	}

	// Check ownership
	if review.UserId != userId {
		return nil, utils.NewError("You can only delete your own review", utils.ErrCodeForbidden)
	}

	courseId := review.CourseId

	// Delete review
	if err := rs.reviewRepo.Delete(reviewId); err != nil {
		return nil, utils.WrapError(err, "Failed to delete review", utils.ErrCodeInternal)
	}

	// Update course rating stats
	if err := rs.reviewRepo.UpdateCourseRatingStats(courseId); err != nil {
		// Log error but don't fail the request
	}

	return &dto.DeleteReviewResponse{
		Message: "Review deleted successfully",
	}, nil
}
