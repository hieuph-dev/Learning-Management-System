package handler

import (
	"lms/src/dto"
	"lms/src/service"
	"lms/src/utils"
	"lms/src/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

// POST /api/v1/payments/create - Create payment
func (ph *PaymentHandler) CreatePayment(ctx *gin.Context) {
	// Get user ID from context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User not authenticated", utils.ErrCodeUnauthorized))
		return
	}

	// Parse request body
	var req dto.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Call service to create payment
	response, err := ph.service.CreatePayment(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/payments/momo/callback - MoMo IPN callback
func (ph *PaymentHandler) MomoCallback(ctx *gin.Context) {
	// Parse callback data
	var callbackData map[string]interface{}
	if err := ctx.ShouldBindJSON(&callbackData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"resultCode": 1,
			"message":    "Invalid request body",
		})
		return
	}

	// Handle callback
	response, err := ph.service.HandleMomoCallback(callbackData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"resultCode": 1,
			"message":    err.Error(),
		})
		return
	}

	// Return response to MoMo
	if response.Success {
		ctx.JSON(http.StatusOK, gin.H{
			"resultCode": 0,
			"message":    "Success",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"resultCode": 1,
			"message":    response.Message,
		})
	}
}

// POST /api/v1/payments/zalopay/callback - ZaloPay callback
func (ph *PaymentHandler) ZaloPayCallback(ctx *gin.Context) {
	// Parse callback data
	var callbackData map[string]interface{}
	if err := ctx.ShouldBindJSON(&callbackData); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"return_code":    0,
			"return_message": "Invalid request body",
		})
		return
	}

	// Handle callback
	response, err := ph.service.HandleZaloPayCallback(callbackData)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"return_code":    0,
			"return_message": err.Error(),
		})
		return
	}

	// Return response to ZaloPay
	if response.Success {
		ctx.JSON(http.StatusOK, gin.H{
			"return_code":    1,
			"return_message": "Success",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"return_code":    0,
			"return_message": response.Message,
		})
	}
}

// GET /api/v1/payments/status - Check payment status
func (ph *PaymentHandler) CheckPaymentStatus(ctx *gin.Context) {
	// Get user ID from context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User not authenticated", utils.ErrCodeUnauthorized))
		return
	}

	// Parse query parameters
	var req dto.CheckPaymentStatusRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Call service to check status
	response, err := ph.service.CheckPaymentStatus(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
