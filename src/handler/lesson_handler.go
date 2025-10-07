package handler

import (
	"lms/src/service"
	"lms/src/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(service service.LessonService) *LessonHandler {
	return &LessonHandler{
		service: service,
	}
}

// GET /api/v1/courses/:course_id/lessons  - Lấy lessons của course (enrolled students only)
func (lh *LessonHandler) GetCourseLessons(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("course_id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course Id format", utils.ErrCodeBadRequest))
	}

	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Gọi service để lấy lessons
	response, err := lh.service.GetCourseLessons(userId.(uint), uint(courseId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/courses/id/:id/lessons/:slug - Lấy lesson detail (enrolled only)
func (lh *LessonHandler) GetLessonDetail(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy lesson slug từ URL parameter
	slug := ctx.Param("slug")
	if slug == "" {
		utils.ResponseError(ctx, utils.NewError("Lesson slug is required", utils.ErrCodeBadRequest))
		return
	}

	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Gọi service để lấy lesson detail
	response, err := lh.service.GetLessonDetail(userId.(uint), uint(courseId), slug)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
