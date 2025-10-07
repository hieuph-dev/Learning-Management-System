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

type EnrollmentHandler struct {
	service service.EnrollmentService
}

func NewEnrollmentHandler(service service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		service: service,
	}
}

// POST /api/v1/courses/:course_id/enroll - Enroll vào course
func (eh *EnrollmentHandler) EnrollCourse(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("course_id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	courseId, err := strconv.ParseInt(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.EnrollCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để enroll
	response, err := eh.service.EnrollCourse(userId.(uint), uint(courseId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, response)
}

// GET /api/v1/course_id/:course_id/check-enrollment - Kiểm tra đã enroll chưa
func (eh *EnrollmentHandler) CheckEnrollment(ctx *gin.Context) {
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
		return
	}

	// Lấy user ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Gọi service để check enrollment
	response, err := eh.service.CheckEnrollment(userId.(uint), uint(courseId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/enrollments/my - Lấy danh sách enrollments của user
func (eh *EnrollmentHandler) GetMyEnrollments(ctx *gin.Context) {
	// Lấy user ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Parse query parameters
	var req dto.GetMyEnrollmentsQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy enrollments
	response, err := eh.service.GetMyEnrollments(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
