package handler

import (
	"lms/src/dto"
	"lms/src/service"
	"lms/src/utils"
	"lms/src/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// POST /api/v1/auth/register
func (ah *AuthHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	// 1. Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// 2. Gọi service để xử lý đăng ký
	createdUser, err := ah.service.Register(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, createdUser)
}

// POST /api/v1/auth/login
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	user, err := ah.service.Login(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, user)
}

// GET /api/v1/auth/profile - PROTECTED
func (ah *AuthHandler) GetProfile(ctx *gin.Context) {
	// Lấy thông tin user từ context (đã được set bởi middleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("user not found in context", utils.ErrCodeUnauthorized))
		return
	}

	profile, err := ah.service.GetProfile(userId.(uint))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, profile)
}

// POST /api/v1/auth/logout - PROTECTED
func (ah *AuthHandler) Logout(ctx *gin.Context) {
	// Trong implementation thực tế, bạn có thể:
	// 1. Blacklist token hiện tại
	// 2. Xóa refresh token khỏi database
	// 3. Log logout event

	utils.ResponseSuccess(ctx, http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// POST /api/v1/auth/refresh - PROTECTED (với refresh token)
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	tokens, err := ah.service.RefreshToken(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, tokens)
}

// POST /api/v1/auth/forgot-password
func (ah *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	response, err := ah.service.ForgotPassword(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

func (ah *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	err := ah.service.ResetPassword(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}
