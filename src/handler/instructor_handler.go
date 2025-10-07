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

type InstructorHandler struct {
	service service.InstructorService
}

func NewInstructorHandler(service service.InstructorService) *InstructorHandler {
	return &InstructorHandler{
		service: service,
	}
}

// GET /api/v1/instructor/courses - Lấy danh sách courses của instructor
func (ih *InstructorHandler) GetInstructorCourses(ctx *gin.Context) {
	// Lấy instructor ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

	// Parse query parameters
	var req dto.GetInstructorCoursesQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy courses
	response, err := ih.service.GetInstructorCourses(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/instructor/courses - Tạo course mới
func (ih *InstructorHandler) CreateCourse(ctx *gin.Context) {
	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để tạo course
	response, err := ih.service.CreateCourse(userId.(uint), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, response)
}

// PUT /api/v1/instructor/courses/:user_id - Cập nhật course
func (ih *InstructorHandler) UpdateCourse(ctx *gin.Context) {
	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

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

	// Bind JSON request
	var req dto.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để cập nhật course
	response, err := ih.service.UpdateCourse(userId.(uint), uint(courseId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// DELETE /api/v1/instructor/courses/:user_id - Xóa course
func (ih *InstructorHandler) DeleteCourse(ctx *gin.Context) {
	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

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

	// Gọi service để xóa course
	response, err := ih.service.DeleteCourse(userId.(uint), uint(courseId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/instructor/courses/:user_id/students - Lấy danh sách students của course
func (ih *InstructorHandler) GetCourseStudents(ctx *gin.Context) {
	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found in context", utils.ErrCodeUnauthorized))
		return
	}

	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("course_id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course ID format", utils.ErrCodeBadRequest))
		return
	}

	// Parse query parameters
	var req dto.GetCourseStudentsQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy students
	response, err := ih.service.GetCourseStudents(userId.(uint), uint(courseId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/instructor/courses/:course_id/lessons - Tạo lesson mới
func (ih *InstructorHandler) CreateLesson(ctx *gin.Context) {
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

	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.CreateLessonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để tạo lesson
	response, err := ih.service.CreateLesson(userId.(uint), uint(courseId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, response)
}

// PUT /api/v1/instructor/courses/:course_id/lessons/:id - Update lesson
func (ih *InstructorHandler) UpdateLesson(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("course_id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy lesson ID từ URL parameter
	lessonIdParam := ctx.Param("id")
	if lessonIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Lesson Id is required", utils.ErrCodeBadRequest))
		return
	}

	lessonId, err := strconv.ParseUint(lessonIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid lesson Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.UpdateLessonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để update lesson
	response, err := ih.service.UpdateLesson(userId.(uint), uint(courseId), uint(lessonId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// DELETE /api/v1/instructor/courses/:course_id/lessons/:id - Delete lesson
func (ih *InstructorHandler) DeleteLesson(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("course_id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy lesson ID từ URL parameter
	lessonIdParam := ctx.Param("id")
	if lessonIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Lesson Id is required", utils.ErrCodeBadRequest))
		return
	}

	lessonId, err := strconv.ParseUint(lessonIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid lesson Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Gọi service để delete lesson
	response, err := ih.service.DeleteLesson(userId.(uint), uint(courseId), uint(lessonId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// PUT /api/v1/instructor/lessons/:id/reorder - Reorder lessons
func (ih *InstructorHandler) ReorderLessons(ctx *gin.Context) {
	// Lấy lesson ID từ URL parameter (có thể là bất kỳ lesson nào trong course)
	lessonIdParam := ctx.Param("id")
	if lessonIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Lesson Id is required", utils.ErrCodeBadRequest))
		return
	}

	lessonId, err := strconv.ParseUint(lessonIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid lesson Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy instructor ID từ context
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Bind JSON request
	var req dto.ReorderLessonsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để reorder lessons
	response, err := ih.service.ReorderLessons(userId.(uint), uint(lessonId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
