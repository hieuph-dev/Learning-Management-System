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

type ProgressHandler struct {
	service service.ProgressService
}

func NewProgressHandler(service service.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		service: service,
	}
}

// GET /api/v1/enrollments/:course_id/progress - Get learning progress
func (ph *ProgressHandler) GetCourseProgress(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("course_id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course ID is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course ID format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Gọi service để lấy progress
	response, err := ph.service.GetCourseProgress(userId.(uint), uint(courseId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/progress/:lesson_id/complete
func (ph *ProgressHandler) CompleteLesson(ctx *gin.Context) {
	// Lấy lesson ID từ URL parameter
	lessonIdParam := ctx.Param("lesson_id")
	if lessonIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Lesson Id is required", utils.ErrCodeBadRequest))
		return
	}

	lessonId, err := strconv.ParseUint(lessonIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid lesson Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy user ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.CompleteLessonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service
	response, err := ph.service.CompleteLesson(userId.(uint), uint(lessonId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// PUT /api/v1/progress/:lesson_id/position
func (ph *ProgressHandler) UpdateLessonPosition(ctx *gin.Context) {
	// Lấy lesson ID từ URL parameter
	lessonIdParam := ctx.Param("lesson_id")
	if lessonIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Lesson Id is required", utils.ErrCodeBadRequest))
		return
	}

	lessonId, err := strconv.ParseUint(lessonIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid lesson Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy user ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.UpdateLessonPositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service
	response, err := ph.service.UpdateLessonPosition(userId.(uint), uint(lessonId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
