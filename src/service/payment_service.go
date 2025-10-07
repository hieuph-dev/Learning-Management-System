package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"lms/src/payment"
	"lms/src/repository"
	"lms/src/utils"
	"time"
)

type paymentService struct {
	orderRepo      repository.OrderRepository
	enrollmentRepo repository.EnrollmentRepository
	couponRepo     repository.CouponRepository
	momoGateway    payment.PaymentGateway
	zaloPayGateway payment.PaymentGateway
}

func NewPaymentService(
	orderRepo repository.OrderRepository,
	enrollmentRepo repository.EnrollmentRepository,
	couponRepo repository.CouponRepository,
	momoGateway payment.PaymentGateway,
	zaloPayGateway payment.PaymentGateway,
) PaymentService {
	return &paymentService{
		orderRepo:      orderRepo,
		enrollmentRepo: enrollmentRepo,
		couponRepo:     couponRepo,
		momoGateway:    momoGateway,
		zaloPayGateway: zaloPayGateway,
	}
}

func (ps *paymentService) CreatePayment(userId uint, req *dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error) {
	// 1. Find order
	order, err := ps.orderRepo.FindById(req.OrderId)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// 2. Verify order belongs to user
	if order.UserId != userId {
		return nil, utils.NewError("Access denied", utils.ErrCodeForbidden)
	}

	// 3. Check order status
	if order.PaymentStatus != "pending" {
		return nil, utils.NewError("Order has already been processed", utils.ErrCodeBadRequest)
	}

	// 4. Check if free order
	if order.FinalPrice == 0 {
		return nil, utils.NewError("This is a free order, no payment required", utils.ErrCodeBadRequest)
	}

	// 5. Select payment gateway
	var gateway payment.PaymentGateway
	switch req.PaymentMethod {
	case dto.PaymentMethodMomo:
		gateway = ps.momoGateway
	case dto.PaymentMethodZaloPay:
		gateway = ps.zaloPayGateway
	default:
		return nil, utils.NewError("Unsupported payment method", utils.ErrCodeBadRequest)
	}

	// 6. Create payment description
	description := fmt.Sprintf("Payment for order %s", order.OrderCode)

	// 7. Create payment with gateway
	response, err := gateway.CreatePayment(order.OrderCode, order.FinalPrice, description)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to create payment", utils.ErrCodeInternal)
	}

	// 8. Update order payment method
	order.PaymentMethod = req.PaymentMethod
	if err := ps.orderRepo.Update(order); err != nil {
		return nil, utils.WrapError(err, "Failed to update order", utils.ErrCodeInternal)
	}

	// 9. Add order info to response
	response.OrderId = order.Id
	response.OrderCode = order.OrderCode
	response.Amount = order.FinalPrice

	return response, nil
}

func (ps *paymentService) HandleMomoCallback(data map[string]interface{}) (*dto.PaymentCallbackResponse, error) {
	// 1. Verify signature
	isValid, err := ps.momoGateway.VerifyCallback(data)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to verify callback", utils.ErrCodeInternal)
	}

	if !isValid {
		return &dto.PaymentCallbackResponse{
			Success: false,
			Message: "Invalid signature",
		}, nil
	}

	// 2. Get order ID from callback
	momoPayment := ps.momoGateway.(*payment.MomoPayment)
	orderCode, err := momoPayment.GetOrderIdFromCallback(data)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get order ID", utils.ErrCodeInternal)
	}

	// 3. Find order
	order, err := ps.orderRepo.FindByOrderCode(orderCode)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// 4. Check payment success
	isSuccess := momoPayment.IsPaymentSuccess(data)

	// 5. Update order status
	if isSuccess {
		if err := ps.completePayment(order); err != nil {
			return nil, err
		}

		return &dto.PaymentCallbackResponse{
			Success: true,
			Message: "Payment successful",
			OrderId: order.Id,
		}, nil
	}

	// Payment failed
	order.PaymentStatus = dto.PaymentStatusFailed
	if err := ps.orderRepo.Update(order); err != nil {
		return nil, utils.WrapError(err, "Failed to update order", utils.ErrCodeInternal)
	}

	return &dto.PaymentCallbackResponse{
		Success: false,
		Message: "Payment failed",
		OrderId: order.Id,
	}, nil
}

func (ps *paymentService) HandleZaloPayCallback(data map[string]interface{}) (*dto.PaymentCallbackResponse, error) {
	// 1. Verify MAC
	isValid, err := ps.zaloPayGateway.VerifyCallback(data)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to verify callback", utils.ErrCodeInternal)
	}

	if !isValid {
		return &dto.PaymentCallbackResponse{
			Success: false,
			Message: "Invalid MAC",
		}, nil
	}

	// 2. Get order ID from callback
	zaloPayment := ps.zaloPayGateway.(*payment.ZaloPayPayment)
	orderCode, err := zaloPayment.GetOrderIdFromCallback(data)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get order ID", utils.ErrCodeInternal)
	}

	// 3. Find order
	order, err := ps.orderRepo.FindByOrderCode(orderCode)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// 4. Check payment success
	isSuccess := zaloPayment.IsPaymentSuccess(data)

	// 5. Update order status
	if isSuccess {
		if err := ps.completePayment(order); err != nil {
			return nil, err
		}

		return &dto.PaymentCallbackResponse{
			Success: true,
			Message: "Payment successful",
			OrderId: order.Id,
		}, nil
	}

	// Payment failed
	order.PaymentStatus = dto.PaymentStatusFailed
	if err := ps.orderRepo.Update(order); err != nil {
		return nil, utils.WrapError(err, "Failed to update order", utils.ErrCodeInternal)
	}

	return &dto.PaymentCallbackResponse{
		Success: false,
		Message: "Payment failed",
		OrderId: order.Id,
	}, nil
}

func (ps *paymentService) CheckPaymentStatus(userId uint, req *dto.CheckPaymentStatusRequest) (*dto.CheckPaymentStatusResponse, error) {
	// 1. Find order
	order, err := ps.orderRepo.FindById(req.OrderId)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// 2. Verify order belongs to user
	if order.UserId != userId {
		return nil, utils.NewError("Access denied", utils.ErrCodeForbidden)
	}

	// 3. If already paid, return current status
	if order.PaymentStatus == dto.PaymentStatusPaid {
		return &dto.CheckPaymentStatusResponse{
			OrderId:       order.Id,
			OrderCode:     order.OrderCode,
			PaymentStatus: order.PaymentStatus,
			PaymentMethod: order.PaymentMethod,
			Amount:        order.FinalPrice,
			PaidAt:        order.PaidAt,
			Message:       "Payment completed",
		}, nil
	}

	// 4. Query payment gateway for status
	var gateway payment.PaymentGateway
	switch order.PaymentMethod {
	case dto.PaymentMethodMomo:
		gateway = ps.momoGateway
	case dto.PaymentMethodZaloPay:
		gateway = ps.zaloPayGateway
	default:
		return &dto.CheckPaymentStatusResponse{
			OrderId:       order.Id,
			OrderCode:     order.OrderCode,
			PaymentStatus: order.PaymentStatus,
			PaymentMethod: order.PaymentMethod,
			Amount:        order.FinalPrice,
			Message:       "Payment pending",
		}, nil
	}

	// 5. Get status from gateway
	status, err := gateway.GetPaymentStatus(order.OrderCode)
	if err != nil {
		// If error, return current status
		return &dto.CheckPaymentStatusResponse{
			OrderId:       order.Id,
			OrderCode:     order.OrderCode,
			PaymentStatus: order.PaymentStatus,
			PaymentMethod: order.PaymentMethod,
			Amount:        order.FinalPrice,
			Message:       "Unable to check payment status",
		}, nil
	}

	// 6. Update order if status changed
	if status != order.PaymentStatus {
		if status == dto.PaymentStatusPaid {
			if err := ps.completePayment(order); err != nil {
				return nil, err
			}
		} else {
			order.PaymentStatus = status
			if err := ps.orderRepo.Update(order); err != nil {
				return nil, utils.WrapError(err, "Failed to update order", utils.ErrCodeInternal)
			}
		}
	}

	return &dto.CheckPaymentStatusResponse{
		OrderId:       order.Id,
		OrderCode:     order.OrderCode,
		PaymentStatus: status,
		PaymentMethod: order.PaymentMethod,
		Amount:        order.FinalPrice,
		PaidAt:        order.PaidAt,
		Message:       ps.getStatusMessage(status),
	}, nil
}

// Helper function to complete payment
func (ps *paymentService) completePayment(order *models.Order) error {
	// 1. Update order status
	now := time.Now()
	order.PaymentStatus = dto.PaymentStatusPaid
	order.PaidAt = &now

	if err := ps.orderRepo.Update(order); err != nil {
		return utils.WrapError(err, "Failed to update order", utils.ErrCodeInternal)
	}

	// 2. Create enrollment if not exists
	if _, exists := ps.enrollmentRepo.CheckEnrollment(order.UserId, order.CourseId); !exists {
		enrollment := &models.Enrollment{
			UserId:             order.UserId,
			CourseId:           order.CourseId,
			EnrolledAt:         now,
			ProgressPercentage: 0,
			Status:             "active",
		}

		if err := ps.enrollmentRepo.Create(enrollment); err != nil {
			return utils.WrapError(err, "Failed to create enrollment", utils.ErrCodeInternal)
		}
	}

	// 3. Update coupon used count
	if order.CouponId != nil {
		if err := ps.couponRepo.IncrementUsedCount(*order.CouponId); err != nil {
			// Log error but don't fail
			fmt.Printf("Failed to increment coupon used count: %v\n", err)
		}
	}

	return nil
}

func (ps *paymentService) getStatusMessage(status string) string {
	switch status {
	case dto.PaymentStatusPaid:
		return "Payment completed successfully"
	case dto.PaymentStatusPending:
		return "Payment is pending"
	case dto.PaymentStatusFailed:
		return "Payment failed"
	case dto.PaymentStatusCancelled:
		return "Payment was cancelled"
	default:
		return "Unknown payment status"
	}
}
