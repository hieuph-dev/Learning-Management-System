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

type AdminHandler struct {
	service      service.AdminService
	orderService service.OrderService
}

func NewAdminHandler(service service.AdminService, orderService service.OrderService) *AdminHandler {
	return &AdminHandler{
		service:      service,
		orderService: orderService,
	}
}

// GET /api/v1/admin/users - Lấy danh sách users (Admin only)
func (ah *AdminHandler) GetUsers(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetUsersQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy danh sách users
	response, err := ah.service.GetUsers(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/admin/users/:id - Lấy thông tin chi tiết user (Admin only)
func (ah *AdminHandler) GetUserById(ctx *gin.Context) {
	// Lấy user ID từ URL parameter
	userIdParam := ctx.Param("id")
	if userIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("User Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid user Id format", utils.ErrCodeBadRequest))
		return
	}

	// Gọi service để lấy thông tin user
	user, err := ah.service.GetUserById(uint(userId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, user)
}

// PUT /api/v1/admin/users/:id - Cập nhật user (Admin only)
func (ah *AdminHandler) UpdateUser(ctx *gin.Context) {
	// Lấy user ID từ URL parameter
	userIdParam := ctx.Param("id")
	if userIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("User Id is required", utils.ErrCodeBadRequest))
	}

	// Convert string to uint
	userId, err := strconv.ParseUint(userIdParam, 18, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid user Id format", utils.ErrCodeBadRequest))
		return
	}

	// Bind JSON request
	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để cập nhật user
	updatedUser, err := ah.service.UpdateUser(uint(userId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, updatedUser)
}

// DELETE /api/v1/admin/users/:id - Xóa user (Admin only)
func (ah *AdminHandler) DeleteUser(ctx *gin.Context) {
	// Lấy user ID từ URL parameter
	userIdParam := ctx.Param("id")
	if userIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("User Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid user Id format", utils.ErrCodeBadRequest))
		return
	}

	// Lấy thông tin admin từ context
	adminId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("Admin information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Không cho phép admin tự xóa chính mình

	if adminId.(uint) == uint(userId) {
		utils.ResponseError(ctx, utils.NewError("Cannot delete your own account", utils.ErrCodeForbidden))
		return
	}

	// Gọi service để xóa user
	response, err := ah.service.DeleteUser(uint(userId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// PUT /api/v1/admin/users/:id/status - Thay đổi trạng thái user (Admin only)
func (ah *AdminHandler) ChangeUserStatus(ctx *gin.Context) {
	// Lấy user ID từ URL parameter
	userIdParam := ctx.Param("id")
	if userIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("User Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	userId, err := strconv.ParseUint(userIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid user Id format", utils.ErrCodeBadRequest))
		return
	}

	// Bind JSON request
	var req dto.ChangeUserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Lấy thông tin admin từ context
	adminId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("Admin information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Không cho phép admin thay đổi trạng thái chính mình
	if adminId.(uint) == uint(userId) {
		utils.ResponseError(ctx, utils.NewError("Cannot change your own account status", utils.ErrCodeForbidden))
		return
	}

	// Gọi service để thay đổi trạng thái
	response, err := ah.service.ChangeUserStatus(uint(userId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/admin/courses - Lấy tất cả courses (Admin)
func (ah *AdminHandler) GetCourses(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetAdminCoursesQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy danh sách courses
	response, err := ah.service.GetCourses(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// PUT /api/v1/admin/courses/:course_id/status - Thay đổi trạng thái course (Admin)
func (ah *AdminHandler) ChangeCourseStatus(ctx *gin.Context) {
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
	var req dto.ChangeCourseStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để thay đổi status
	response, err := ah.service.ChangeCourseStatus(uint(courseId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/admin/orders - Lấy tất cả orders (Admin)
func (ah *AdminHandler) GetAllOrders(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetAdminOrdersQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy orders
	response, err := ah.orderService.GetAllOrders(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// PUT /api/v1/admin/orders/:id/status - Cập nhật trạng thái order (Admin)
func (ah *AdminHandler) UpdateOrderStatus(ctx *gin.Context) {
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

	// Parse request body
	var req dto.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để update order status
	response, err := ah.orderService.UpdateOrderStatus(uint(orderId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
