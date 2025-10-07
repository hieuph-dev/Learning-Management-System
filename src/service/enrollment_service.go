package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"math"
	"time"

	"github.com/google/uuid"
)

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	orderRepo      repository.OrderRepository
	courseRepo     repository.CourseRepository
	couponRepo     repository.CouponRepository
	progressRepo   repository.ProgressRepository // Thêm để đếm completed lessons
}

func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	orderRepo repository.OrderRepository,
	courseRepo repository.CourseRepository,
	couponRepo repository.CouponRepository,
	progressRepo repository.ProgressRepository,
) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		orderRepo:      orderRepo,
		courseRepo:     courseRepo,
		couponRepo:     couponRepo,
		progressRepo:   progressRepo,
	}
}

func (es *enrollmentService) EnrollCourse(userId, courseId uint, req *dto.EnrollCourseRequest) (*dto.EnrollCourseResponse, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := es.courseRepo.FindById(courseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra course status
	if course.Status != "published" {
		return nil, utils.NewError("Course is not available for enrollment", utils.ErrCodeBadRequest)
	}

	// 3. Kiểm tra user đã enroll chưa
	if existingEnrollment, exists := es.enrollmentRepo.CheckEnrollment(userId, courseId); exists {
		if existingEnrollment.Status == "active" {
			return nil, utils.NewError("You are already enrolled in this course", utils.ErrCodeConflict)
		}
	}

	// 4. Tính toán giá
	originalPrice := course.Price
	if course.DiscountPrice != nil && *course.DiscountPrice < originalPrice {
		originalPrice = *course.DiscountPrice
	}

	discountAmount := 0.0
	var couponId *uint

	// 5. Áp dụng coupon nếu có
	if req.CouponCode != "" {
		coupon, err := es.couponRepo.FindByCode(req.CouponCode)
		if err != nil {
			return nil, utils.NewError("Invalid coupon code", utils.ErrCodeBadRequest)
		}

		if !es.couponRepo.IsValidCoupon(coupon) {
			return nil, utils.NewError("Coupon is expired or invalid", utils.ErrCodeBadRequest)
		}

		// Check minimum order amount
		if originalPrice < coupon.MinOrderAmount {
			return nil, utils.NewError(
				fmt.Sprintf("Minimum order amount for this coupon is %2.f", coupon.MinOrderAmount),
				utils.ErrCodeBadRequest,
			)
		}

		// Calculate discount
		if coupon.DiscountType == "percentage" {
			discountAmount = originalPrice * (coupon.DiscountValue / 100)
		} else if coupon.DiscountType == "fixed" {
			discountAmount = coupon.DiscountValue
		}

		// Apply max discount if set
		if coupon.MaxDiscountAmount != nil && discountAmount > *coupon.MaxDiscountAmount {
			discountAmount = *coupon.MaxDiscountAmount
		}

		couponId = &coupon.Id
	}

	// 6. Tính final price
	finalPrice := originalPrice - discountAmount
	if finalPrice < 0 {
		finalPrice = 0
	}

	// 7. Tạo order code
	orderCode := fmt.Sprintf("ORD-%s-%d", uuid.New().String()[:8], time.Now().Unix())

	// 8. Tạo order
	order := &models.Order{
		UserId:         userId,
		CourseId:       courseId,
		OrderCode:      orderCode,
		OriginalPrice:  originalPrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		CouponId:       couponId,
		PaymentMethod:  req.PaymentMethod,
		PaymentStatus:  "pending",
	}

	if err := es.orderRepo.Create(order); err != nil {
		return nil, utils.WrapError(err, "Failed to create order", utils.ErrCodeInternal)
	}

	// 9. Nếu course free (finalPrice = 0), tự động approve
	if finalPrice == 0 {
		order.PaymentStatus = "paid"
		now := time.Now()
		order.PaidAt = &now

		if err := es.orderRepo.UpdatePaymentStatus(order.Id, "paid"); err != nil {
			return nil, utils.WrapError(err, "Failed to update payment status", utils.ErrCodeInternal)
		}
	}

	// 10. Tạo enrollment
	enrollment := &models.Enrollment{
		UserId:             userId,
		CourseId:           courseId,
		EnrolledAt:         time.Now(),
		ProgressPercentage: 0,
		Status:             "active",
	}

	if err := es.enrollmentRepo.Create(enrollment); err != nil {
		return nil, utils.WrapError(err, "Failed to create enrollment", utils.ErrCodeInternal)
	}

	// 11. Update coupon used count nếu có
	if couponId != nil {
		if err := es.couponRepo.IncrementUsedCount(*couponId); err != nil {
			// Log error but don't fail
			fmt.Printf("Failed to increment coupon used count: %v\n", err)
		}
	}

	// 12. Update course enrolled count
	// TODO: Implement UpdateEnrolledCount in CourseRepository

	return &dto.EnrollCourseResponse{
		EnrollmentId:   enrollment.Id,
		OrderId:        order.Id,
		OrderCode:      order.OrderCode,
		CourseId:       course.Id,
		CourseTitle:    course.Title,
		OriginalPrice:  originalPrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		PaymentMethod:  order.PaymentMethod,
		PaymentStatus:  order.PaymentStatus,
		EnrolledAt:     enrollment.EnrolledAt,
		Message:        getEnrollmentMessage(finalPrice, order.PaymentStatus),
	}, nil
}

func (es *enrollmentService) CheckEnrollment(userId, courseId uint) (*dto.CheckEnrollmentResponse, error) {
	enrollment, exists := es.enrollmentRepo.CheckEnrollment(userId, courseId)

	if !exists {
		return &dto.CheckEnrollmentResponse{
			IsEnrolled: false,
		}, nil
	}

	return &dto.CheckEnrollmentResponse{
		IsEnrolled:      true,
		EnrollmentId:    enrollment.Id,
		EnrolledAt:      &enrollment.EnrolledAt,
		ProgressPercent: enrollment.ProgressPercentage,
	}, nil
}

func (es *enrollmentService) GetMyEnrollments(userId uint, req *dto.GetMyEnrollmentsQueryRequest) (*dto.GetMyEnrollmentsResponse, error) {
	// Set defaults
	page := 1
	limit := 10

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})
	if req.Status != "" {
		filters["status"] = req.Status
	}

	// Get enrollments
	enrollments, total, err := es.enrollmentRepo.GetUserEnrollments(userId, offset, limit, filters)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get enrollments", utils.ErrCodeInternal)
	}

	// Convert to DTO with course details
	enrollmentItems := make([]dto.EnrollmentItem, len(enrollments))
	for i, enrollment := range enrollments {
		// Load course info
		course, err := es.courseRepo.FindById(enrollment.CourseId)
		if err != nil {
			// Nếu course không tìm thấy, skip hoặc trả về thông tin cơ bản
			enrollmentItems[i] = dto.EnrollmentItem{
				Id:                 enrollment.Id,
				CourseId:           enrollment.CourseId,
				CourseTitle:        "Course not found",
				CourseThumbnail:    "",
				InstructorName:     "",
				EnrolledAt:         enrollment.EnrolledAt,
				ProgressPercentage: enrollment.ProgressPercentage,
				Status:             enrollment.Status,
				TotalLessons:       0,
				CompletedLessons:   0,
			}
			continue
		}

		// Count completed lessons
		completedLessons := 0
		if es.progressRepo != nil {
			completedLessons, _ = es.progressRepo.CountCompletedLessons(userId, enrollment.CourseId)
		}

		enrollmentItems[i] = dto.EnrollmentItem{
			Id:                 enrollment.Id,
			CourseId:           enrollment.CourseId,
			CourseTitle:        course.Title,
			CourseThumbnail:    course.ThumbnailURL,
			InstructorName:     course.Instructor.FullName,
			EnrolledAt:         enrollment.EnrolledAt,
			ProgressPercentage: enrollment.ProgressPercentage,
			Status:             enrollment.Status,
			TotalLessons:       course.TotalLessons,
			CompletedLessons:   completedLessons,
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

	return &dto.GetMyEnrollmentsResponse{
		Enrollments: enrollmentItems,
		Pagination:  pagination,
	}, nil
}

func getEnrollmentMessage(finalPrice float64, paymentStatus string) string {
	if finalPrice == 0 {
		return "Congratulations! You have successfully enrolled in this free course"
	}

	if paymentStatus == "paid" {
		return "Payment successfull! You have been enrolled in the course"
	}

	return "Enrollment pending. Please complete your payment to access the course"
}
