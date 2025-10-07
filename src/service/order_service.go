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

type orderService struct {
	orderRepo      repository.OrderRepository
	courseRepo     repository.CourseRepository
	couponRepo     repository.CouponRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	courseRepo repository.CourseRepository,
	couponRepo repository.CouponRepository,
	enrollmentRepo repository.EnrollmentRepository,

) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
		couponRepo:     couponRepo,
	}
}

func (os *orderService) CreateOrder(userId uint, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := os.courseRepo.FindById(req.CourseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra course status
	if course.Status != "published" {
		return nil, utils.NewError("Course is not available for purchase", utils.ErrCodeBadRequest)
	}

	// 3. Kiểm tra user đã mua course chưa
	if existingEnrollment, exists := os.enrollmentRepo.CheckEnrollment(userId, req.CourseId); exists {
		if existingEnrollment.Status == "active" {
			return nil, utils.NewError("You already own this course", utils.ErrCodeConflict)
		}
	}

	// 4. Kiểm tra đã có order pending chưa
	existingOrder, err := os.orderRepo.FindPendingOrderByUserAndCourse(userId, req.CourseId)
	if err == nil && existingOrder != nil {
		return nil, utils.NewError("You already have a pending order for this course. Please complete or cancel it first", utils.ErrCodeConflict)
	}

	// 5. Tính giá gốc
	originalPrice := course.Price
	if course.DiscountPrice != nil && *course.DiscountPrice < originalPrice {
		originalPrice = *course.DiscountPrice
	}

	discountAmount := 0.0
	var couponId *uint
	var appliedCouponCode string

	// 6. Áp dụng coupon nếu có
	if req.CouponCode != "" {
		coupon, err := os.couponRepo.FindByCode(req.CouponCode)
		if err != nil {
			return nil, utils.NewError("Invalid coupon code", utils.ErrCodeBadRequest)
		}

		if !os.couponRepo.IsValidCoupon(coupon) {
			return nil, utils.NewError("Coupon is expired or not available", utils.ErrCodeBadRequest)
		}

		// Kiểm tra minimum order amount
		if originalPrice < coupon.MinOrderAmount {
			return nil, utils.NewError(
				fmt.Sprintf("Minimum order amount for this coupon is %2.f", coupon.MinOrderAmount),
				utils.ErrCodeBadRequest,
			)
		}

		// Tính discount
		if coupon.DiscountType == "percentage" {
			discountAmount = originalPrice * (coupon.DiscountValue / 100)
		} else if coupon.DiscountType == "fixed" {
			discountAmount = coupon.DiscountValue
		}

		// Apply max discount nếu có
		if coupon.MaxDiscountAmount != nil && discountAmount > *coupon.MaxDiscountAmount {
			discountAmount = *coupon.MaxDiscountAmount
		}

		couponId = &coupon.Id
		appliedCouponCode = coupon.Code
	}

	// 7. Tính final price
	finalPrice := originalPrice - discountAmount
	if finalPrice < 0 {
		finalPrice = 0
	}

	// 8. Tạo order code
	orderCode := fmt.Sprintf("ORD-%s-%d", uuid.New().String()[:8], time.Now().Unix())

	// 9. Tạo order
	order := &models.Order{
		UserId:         userId,
		CourseId:       req.CourseId,
		OrderCode:      orderCode,
		OriginalPrice:  originalPrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		CouponId:       couponId,
		PaymentStatus:  "pending",
	}

	if err := os.orderRepo.Create(order); err != nil {
		return nil, utils.WrapError(err, "Failed to create order", utils.ErrCodeInternal)
	}

	// 10. Nếu free course, tự động approve và tạo enrollment
	message := "Order created successfully. Please proceed to payment"
	if finalPrice == 0 {
		if err := os.completeOrder(order, "free"); err != nil {
			return nil, err
		}
		message = "Congratulations! You have successfully enrolled in this free course"
	}

	return &dto.CreateOrderResponse{
		OrderId:        order.Id,
		OrderCode:      order.OrderCode,
		CourseId:       course.Id,
		CourseTitle:    course.Title,
		OriginalPrice:  originalPrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		CouponCode:     appliedCouponCode,
		PaymentStatus:  order.PaymentStatus,
		CreatedAt:      order.CreatedAt,
		Message:        message,
	}, nil
}

func (os *orderService) GetOrderHistory(userId uint, req *dto.GetOrderHistoryQueryRequest) (*dto.GetOrderHistoryResponse, error) {
	// Set defaults
	page := 1
	limit := 10
	sortBy := "desc"

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}
	if req.SortBy != "" {
		sortBy = req.SortBy
	}

	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})
	if req.PaymentStatus != "" {
		filters["payment_status"] = req.PaymentStatus
	}

	// Get orders
	orders, total, err := os.orderRepo.GetUsersOrders(userId, offset, limit, filters, "created_at", sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get order history", utils.ErrCodeInternal)
	}

	// Convert to DTO
	orderItems := make([]dto.OrderHistoryItem, len(orders))
	for i, order := range orders {
		// Load course info
		course, err := os.courseRepo.FindById(order.CourseId)
		if err != nil {
			orderItems[i] = dto.OrderHistoryItem{
				Id:              order.Id,
				OrderCode:       order.OrderCode,
				CourseId:        order.CourseId,
				CourseTitle:     "Course not found",
				CourseThumbnail: "",
				OriginalPrice:   order.OriginalPrice,
				DiscountAmount:  order.DiscountAmount,
				FinalPrice:      order.FinalPrice,
				PaymentStatus:   order.PaymentStatus,
				PaidAt:          order.PaidAt,
				CreatedAt:       order.CreatedAt,
			}
			continue
		}

		orderItems[i] = dto.OrderHistoryItem{
			Id:              order.Id,
			OrderCode:       order.OrderCode,
			CourseId:        order.CourseId,
			CourseTitle:     course.Title,
			CourseThumbnail: course.ThumbnailURL,
			OriginalPrice:   order.OriginalPrice,
			DiscountAmount:  order.DiscountAmount,
			FinalPrice:      order.FinalPrice,
			PaymentStatus:   order.PaymentStatus,
			PaidAt:          order.PaidAt,
			CreatedAt:       order.CreatedAt,
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

	return &dto.GetOrderHistoryResponse{
		Orders:     orderItems,
		Pagination: pagination,
	}, nil
}

// Helper function to complete order and create enrollment
func (os *orderService) completeOrder(order *models.Order, paymentMethod string) error {
	now := time.Now()
	order.PaymentStatus = "paid"
	order.PaymentMethod = paymentMethod
	order.PaidAt = &now

	// Update order
	if err := os.orderRepo.Update(order); err != nil {
		return utils.WrapError(err, "Failed to update order", utils.ErrCodeInternal)
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		UserId:             order.UserId,
		CourseId:           order.CourseId,
		EnrolledAt:         now,
		ProgressPercentage: 0,
		Status:             "active",
	}

	if err := os.enrollmentRepo.Create(enrollment); err != nil {
		return utils.WrapError(err, "Failed to create enrollment", utils.ErrCodeInternal)
	}

	// Update coupon used count nếu có
	if order.CouponId != nil {
		os.couponRepo.IncrementUsedCount(*order.CouponId)
	}

	return nil
}

func (os *orderService) GetOrderDetail(userId uint, orderId uint) (*dto.OrderDetailResponse, error) {
	// Tìm order
	order, err := os.orderRepo.FindById(orderId)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// Kiểm tra order có thuộc về user không
	if order.UserId != userId {
		return nil, utils.NewError("Access denied", utils.ErrCodeForbidden)
	}

	// Load course info
	course, err := os.courseRepo.FindById(order.CourseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// Get coupon code nếu có
	couponCode := ""
	if order.CouponId != nil {
		if coupon, err := os.couponRepo.FindById(*order.CouponId); err == nil {
			couponCode = coupon.Code
		}
	}

	return &dto.OrderDetailResponse{
		Id:              order.Id,
		OrderCode:       order.OrderCode,
		UserId:          order.UserId,
		CourseId:        order.CourseId,
		CourseTitle:     course.Title,
		CourseThumbnail: course.ThumbnailURL,
		InstructorName:  course.Instructor.FullName,
		OriginalPrice:   order.OriginalPrice,
		DiscountAmount:  order.DiscountAmount,
		FinalPrice:      order.FinalPrice,
		CouponCode:      couponCode,
		PaymentStatus:   order.PaymentStatus,
		PaidAt:          order.PaidAt,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}, nil
}

func (os *orderService) PayOrder(userId uint, orderId uint, req *dto.PayOrderRequest) (*dto.PayOrderResponse, error) {
	// 1. Tìm order
	order, err := os.orderRepo.FindById(orderId)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra order có thuộc về user không
	if order.UserId != userId {
		return nil, utils.NewError("Access denied", utils.ErrCodeForbidden)
	}

	// 3. Kiểm tra order status
	if order.PaymentStatus != "pending" {
		return nil, utils.NewError("Order has already been processed", utils.ErrCodeBadRequest)
	}

	// 4. Kiểm tra nếu là free course
	if order.FinalPrice == 0 {
		return nil, utils.NewError("This is a free order, no payment required", utils.ErrCodeBadRequest)
	}

	// 5. Simulate payment processing
	// TODO: Integrate with real payment gateway
	time.Sleep(1 * time.Second) // Simulate payment processing

	// 6. Complete order
	order.PaymentMethod = req.PaymentMethod
	if err := os.completeOrder(order, req.PaymentMethod); err != nil {
		return nil, err
	}

	return &dto.PayOrderResponse{
		OrderId:       order.Id,
		OrderCode:     order.OrderCode,
		PaymentStatus: "paid",
		PaymentMethod: req.PaymentMethod,
		PaidAt:        *order.PaidAt,
		Message:       "Payment successful! You have been enrolled in the course",
	}, nil
}

func (os *orderService) GetAllOrders(req *dto.GetAdminOrdersQueryRequest) (*dto.GetAdminOrdersResponse, error) {
	// Set defaults
	page := 1
	limit := 20
	orderBy := "created_at"
	sortBy := "desc"

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}
	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}
	if req.SortBy != "" {
		sortBy = req.SortBy
	}

	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})
	if req.UserId != nil {
		filters["user_id"] = *req.UserId
	}
	if req.CourseId != nil {
		filters["course_id"] = *req.CourseId
	}
	if req.PaymentStatus != "" {
		filters["payment_status"] = req.PaymentStatus
	}
	if req.PaymentMethod != "" {
		filters["payment_method"] = req.PaymentMethod
	}
	if req.Search != "" {
		filters["search"] = req.Search
	}
	if req.DateFrom != "" {
		if dateFrom, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
			filters["date_from"] = dateFrom
		}
	}
	if req.DateTo != "" {
		if dateTo, err := time.Parse("2006-01-02", req.DateTo); err == nil {
			// Set to end of day
			dateTo = dateTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters["date_to"] = dateTo
		}
	}

	// Get orders
	orders, total, err := os.orderRepo.GetAllOrders(offset, limit, filters, orderBy, sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get orders", utils.ErrCodeInternal)
	}

	// Get statistics
	statistics, err := os.orderRepo.GetOrderStatistics(filters)
	if err != nil {
		// Log error but don't fail the request
		statistics = &dto.OrderStatistics{}
	}

	// Convert to DTO
	orderItems := make([]dto.AdminOrderItem, len(orders))
	for i, order := range orders {
		// Get user info
		username := ""
		userEmail := ""
		if order.User.Id != 0 {
			username = order.User.Username
			userEmail = order.User.Email
		}

		// Get course info
		courseTitle := "Course not found"
		courseThumbnail := ""
		instructorName := ""
		if order.Course.Id != 0 {
			courseTitle = order.Course.Title
			courseThumbnail = order.Course.ThumbnailURL
			if order.Course.Instructor.Id != 0 {
				instructorName = order.Course.Instructor.FullName
			}
		}

		// Get coupon code
		couponCode := ""
		if order.CouponId != nil {
			if coupon, err := os.couponRepo.FindById(*order.CouponId); err == nil {
				couponCode = coupon.Code
			}
		}

		orderItems[i] = dto.AdminOrderItem{
			Id:              order.Id,
			OrderCode:       order.OrderCode,
			UserId:          order.UserId,
			Username:        username,
			UserEmail:       userEmail,
			CourseId:        order.CourseId,
			CourseTitle:     courseTitle,
			CourseThumbnail: courseThumbnail,
			InstructorName:  instructorName,
			OriginalPrice:   order.OriginalPrice,
			DiscountAmount:  order.DiscountAmount,
			FinalPrice:      order.FinalPrice,
			CouponCode:      couponCode,
			PaymentMethod:   order.PaymentMethod,
			PaymentStatus:   order.PaymentStatus,
			PaidAt:          order.PaidAt,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
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

	return &dto.GetAdminOrdersResponse{
		Orders:     orderItems,
		Pagination: pagination,
		Statistics: *statistics,
	}, nil
}

func (os *orderService) UpdateOrderStatus(orderId uint, req *dto.UpdateOrderStatusRequest) (*dto.UpdateOrderStatusResponse, error) {
	// 1. Find order
	order, err := os.orderRepo.FindById(orderId)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// 2. Validate status change
	if order.PaymentStatus == req.Status {
		return nil, utils.NewError(
			fmt.Sprintf("Order already has status: %s", req.Status),
			utils.ErrCodeBadRequest,
		)
	}

	// 3. Business logic validation
	// Không cho phép thay đổi từ paid sang pending
	if order.PaymentStatus == "paid" && req.Status == "pending" {
		return nil, utils.NewError("Cannot change paid order back to pending", utils.ErrCodeBadRequest)
	}

	// Không cho phép thay đổi từ cancelled/failed sang paid
	if (order.PaymentStatus == "cancelled" || order.PaymentStatus == "failed") && req.Status == "paid" {
		return nil, utils.NewError(
			fmt.Sprintf("Cannot change %s order to paid. Pleasse create a new order", order.PaymentStatus),
			utils.ErrCodeBadRequest,
		)
	}

	// 4. Handle status change to 'paid'
	if req.Status == "paid" {
		// Create enrollment if not exists
		if _, exists := os.enrollmentRepo.CheckEnrollment(order.UserId, order.CourseId); !exists {
			enrollment := &models.Enrollment{
				UserId:             order.UserId,
				CourseId:           order.CourseId,
				EnrolledAt:         time.Now(),
				ProgressPercentage: 0,
				Status:             "active",
			}

			if err := os.enrollmentRepo.Create(enrollment); err != nil {
				return nil, utils.WrapError(err, "Failed to create enrollment", utils.ErrCodeInternal)
			}
		}

		// Increment coupon used count
		if order.CouponId != nil {
			os.couponRepo.IncrementUsedCount(*order.CouponId)
		}
	}

	// 5. Update order status
	if err := os.orderRepo.UpdateOrderStatus(orderId, req.Status); err != nil {
		return nil, utils.WrapError(err, "Failed to update order status", utils.ErrCodeInternal)
	}

	// 6. Get updated order
	updatedOrder, err := os.orderRepo.FindById(orderId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get updated order", utils.ErrCodeInternal)
	}

	// 7. Build message
	message := fmt.Sprintf("Order status changed from '%s' to '%s'", order.PaymentStatus, req.Status)
	if req.Reason != "" {
		message += fmt.Sprintf(". Reason: %s", req.Reason)
	}

	return &dto.UpdateOrderStatusResponse{
		Id:            updatedOrder.Id,
		OrderCode:     updatedOrder.OrderCode,
		PaymentStatus: updatedOrder.PaymentStatus,
		UpdatedAt:     updatedOrder.UpdatedAt,
		Message:       message,
	}, nil
}
