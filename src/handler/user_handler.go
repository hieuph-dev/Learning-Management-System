package handler

import (
	"lms/src/dto"
	"lms/src/service"
	"lms/src/utils"
	"lms/src/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GET /api/v1/users/profile - Lấy thông tin profile (Auth required)
func (uh *UserHandler) GetProfile(ctx *gin.Context) {
	// Lấy thông tin user từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

	// Gọi service để lấy profile
	profile, err := uh.service.GetProfile(userId.(uint))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, profile)
}

// PUT /api/v1/users/profile - Cập nhật profile (Auth required)
func (uh *UserHandler) UpdateProfile(ctx *gin.Context) {
	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
	}

	// Gọi service để cập nhật profile
	updatedProfile, err := uh.service.UpdateProfile(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, updatedProfile)
}

func (uh *UserHandler) ChangePassword(ctx *gin.Context) {
	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User infomation not found in context", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để đổi mật khẩu
	response, err := uh.service.ChangePassword(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/users/upload-avatar - Upload avatar (Auth required)
func (uh *UserHandler) UploadAvatar(ctx *gin.Context) {
	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

	// Lấy file từ form data
	file, err := ctx.FormFile("avatar")
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Avatar file is required", utils.ErrCodeBadRequest))
		return
	}

	// Gọi service để xử lý upload
	response, err := uh.service.UploadAvatar(userId.(uint), file)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
