package handler

import (
	"lms/src/dto"
	"lms/src/service"
	"lms/src/utils"
	"lms/src/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService  service.OrderService
	couponService service.CouponService
}

func NewOrderHandler(orderService service.OrderService, couponService service.CouponService) *OrderHandler {
	return &OrderHandler{
		orderService:  orderService,
		couponService: couponService,
	}
}

// POST /api/v1/orders - Create order
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	// Lấy userId từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User not authenticated", utils.ErrCodeUnauthorized))
	}

	// Parse request body
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để tạo order
	response, err := oh.orderService.CreateOrder(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, response)
}

// GET /api/v1/orders - Get order history
func (oh *OrderHandler) GetOrderHistory(ctx *gin.Context) {
	// Lấy userId từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User not authenticated", utils.ErrCodeUnauthorized))
		return
	}

	// Parse query parameters
	var req dto.GetOrderHistoryQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy order history
	response, err := oh.orderService.GetOrderHistory(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)

}

// GET /api/v1/orders/:id - Get order details
func (oh *OrderHandler) GetOrderDetail(ctx *gin.Context) {
	// Lấy userId từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User not authenticated", utils.ErrCodeUnauthorized))
		return
	}

	// Lấy order ID từ URL parameter
	orderIdParam := ctx.Param("id")
	if orderIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Order Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	orderId, err := strconv.ParseUint(orderIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid order Id format", utils.ErrCodeBadRequest))
		return
	}

	// Gọi service để lấy order detail
	response, err := oh.orderService.GetOrderDetail(userId.(uint), uint(orderId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/orders/:id/pay - Pay order (simulated)
func (oh *OrderHandler) PayOrder(ctx *gin.Context) {
	// Lấy userId từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User not authenticated", utils.ErrCodeUnauthorized))
		return
	}

	// Lấy order ID từ URL parameter
	orderIdParam := ctx.Param("id")
	if orderIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Order Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	orderId, err := strconv.ParseUint(orderIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid order ID format", utils.ErrCodeBadRequest))
		return
	}

	// Parse request body
	var req dto.PayOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để thanh toán order
	response, err := oh.orderService.PayOrder(userId.(uint), uint(orderId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/coupons/validate - Validate coupon
func (oh *OrderHandler) ValidateCoupon(ctx *gin.Context) {
	// Parse request body
	var req dto.ValidateCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để validate coupon
	response, err := oh.couponService.ValidateCoupon(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
