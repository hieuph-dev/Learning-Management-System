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

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

// GET /api/v1/categories - Lấy danh sách categories (Public)
func (ch *CategoryHandler) GetCategories(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetCategoriesQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy categories
	response, err := ch.service.GetCategories(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/categories/:id - Lấy thông tin category (Public)
func (ch *CategoryHandler) GetCategoryById(ctx *gin.Context) {
	// Lấy category ID từ URL parameter
	categoryIdParam := ctx.Param("id")
	if categoryIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Category Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	categoryId, err := strconv.ParseUint(categoryIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid category Id format", utils.ErrCodeBadRequest))
		return
	}

	// Gọi service để lấy thông tin category
	category, err := ch.service.GetCategoryById(uint(categoryId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, category)
}

// POST /api/v1/admin/categories - Tạo category (Admin only)
func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	// Bind JSON request
	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để tạo category
	createdCategory, err := ch.service.CreateCategory(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, createdCategory)
}

// PUT /api/v1/admin/categories/:id - Cập nhật category (Admin only)
func (ch *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	// Lấy category ID từ URL parameter
	categoryIdParam := ctx.Param("id")
	if categoryIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Category Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	categoryId, err := strconv.ParseUint(categoryIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid category Id format", utils.ErrCodeBadRequest))
		return
	}

	// Bind JSON request
	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để cập nhật category
	updatedCategory, err := ch.service.UpdateCategory(uint(categoryId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, updatedCategory)
}

// DELETE /api/v1/admin/categories/:id - Xóa category (Admin only)
func (ch *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	// Lấy category ID từ URL parameter
	categoryIdParam := ctx.Param("id")
	if categoryIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Category Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	categoryId, err := strconv.ParseUint(categoryIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid category Id format", utils.ErrCodeBadRequest))
		return
	}

	// Gọi service để xóa category
	response, err := ch.service.DeleteCategory(uint(categoryId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
