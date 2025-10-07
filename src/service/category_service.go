package service

import (
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"strings"
	"time"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (cs *categoryService) GetCategories(req *dto.GetCategoriesQueryRequest) (*dto.GetCategoriesResponse, error) {
	// Prepare filters
	filters := make(map[string]interface{})

	if req.ParentId != nil {
		filters["parent_id"] = *req.ParentId
	}

	if req.IsActive != nil {
		filters["is_active"] = *req.IsActive
	}

	if req.Search != "" {
		filters["search"] = utils.NormalizeString(req.Search)
	}

	// Get categories
	categories, total, err := cs.categoryRepo.GetCategories(filters)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get categories", utils.ErrCodeInternal)
	}

	// Convert to DTO
	categoryItems := make([]dto.CategoryItem, len(categories))
	for i, category := range categories {
		categoryItems[i] = dto.CategoryItem{
			Id:          category.Id,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			ImageURL:    category.ImageURL,
			ParentId:    category.ParentId,
			SortOrder:   category.SortOrder,
			IsActive:    category.IsActive,
			CreatedAt:   category.CreatedAt,
		}
	}

	return &dto.GetCategoriesResponse{
		Categories: categoryItems,
		Total:      total,
	}, nil
}

func (cs *categoryService) GetCategoryById(categoryId uint) (*dto.CategoryDetail, error) {
	// Tìm category theo ID
	category, err := cs.categoryRepo.FindById(categoryId)
	if err != nil {
		return nil, utils.NewError("Category not found", utils.ErrCodeNotFound)
	}

	// Convert sang DTO
	return &dto.CategoryDetail{
		Id:          category.Id,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageURL:    category.ImageURL,
		ParentId:    category.ParentId,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (cs *categoryService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CreateCategoryResponse, error) {
	// 1. Kiểm tra slug có tồn tại chưa
	req.Slug = utils.NormalizeString(req.Slug)
	if _, exists := cs.categoryRepo.FindBySlug(req.Slug); exists {
		return nil, utils.NewError("Slug already exists", utils.ErrCodeConflict)
	}

	// 2. Kiểm tra parent category nếu có
	if req.ParentId != nil {
		_, err := cs.categoryRepo.FindById(*req.ParentId)
		if err != nil {
			return nil, utils.NewError("Parent category not found", utils.ErrCodeNotFound)
		}
	}

	// 3. Tạo category mới
	category := models.Category{
		Name:        strings.TrimSpace(req.Name),
		Slug:        req.Slug,
		Description: strings.TrimSpace(req.Description),
		ImageURL:    strings.TrimSpace(req.ImageURL),
		ParentId:    req.ParentId,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}

	// 4. Lưu vào database
	if err := cs.categoryRepo.Create(&category); err != nil {
		return nil, utils.WrapError(err, "Failed to create category", utils.ErrCodeInternal)
	}

	return &dto.CreateCategoryResponse{
		Id:          category.Id,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageURL:    category.ImageURL,
		ParentId:    category.ParentId,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
	}, nil
}

func (cs *categoryService) UpdateCategory(categoryId uint, req *dto.UpdateCategoryRequest) (*dto.UpdateCategoryResponse, error) {
	// 1. Kiểm tra category có tồn tại không
	existingCategory, err := cs.categoryRepo.FindById(categoryId)
	if err != nil {
		return nil, utils.NewError("Category not found", utils.ErrCodeNotFound)
	}

	if !existingCategory.IsActive {
		return nil, utils.NewError("Category is not active", utils.ErrCodeForbidden)
	}

	// 2. Chuẩn bị dữ liệu cập nhật
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}

	if req.Slug != "" {
		normalizedSlug := utils.NormalizeString(req.Slug)
		// Kiểm tra slug có trùng với category khác không
		if _, exists := cs.categoryRepo.FindBySlugExcept(normalizedSlug, categoryId); exists {
			return nil, utils.NewError("Slug already exists", utils.ErrCodeConflict)
		}
		updates["slug"] = normalizedSlug
	}

	if req.Description != "" {
		updates["description"] = strings.TrimSpace(req.Description)
	}

	// ParentId có thể là nil nên xử lý riêng
	if req.ParentId != nil {
		// Kiểm tra parent category có tồn tại không
		_, err := cs.categoryRepo.FindById(*req.ParentId)
		if err != nil {
			return nil, utils.NewError("Parent category not found", utils.ErrCodeNotFound)
		}

		// Không cho phép set parent là chính nó
		if *req.ParentId == categoryId {
			return nil, utils.NewError("Category cannot be parent of itself", utils.ErrCodeBadRequest)
		}
		updates["parent_id"] = *req.ParentId
	}

	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	updates["updated_at"] = time.Now()

	// 3. Cập nhật category
	if len(updates) > 1 { // > 1 vì luôn ó updated_at
		if err := cs.categoryRepo.Update(categoryId, updates); err != nil {
			return nil, utils.WrapError(err, "Failed to update category", utils.ErrCodeInternal)
		}
	}

	// 4. Lấy thông tin category đã cập nhật
	updatedCategory, err := cs.categoryRepo.FindById(categoryId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get updated category", utils.ErrCodeInternal)
	}

	return &dto.UpdateCategoryResponse{
		Id:          updatedCategory.Id,
		Name:        updatedCategory.Name,
		Slug:        updatedCategory.Slug,
		Description: updatedCategory.Description,
		ImageURL:    updatedCategory.ImageURL,
		ParentId:    updatedCategory.ParentId,
		SortOrder:   updatedCategory.SortOrder,
		IsActive:    updatedCategory.IsActive,
		CreatedAt:   updatedCategory.CreatedAt,
		UpdatedAt:   updatedCategory.UpdatedAt,
	}, nil
}

func (cs *categoryService) DeleteCategory(categoryId uint) (*dto.DeleteCategoryResponse, error) {
	// 1. Kiểm tra category có tồn tại không
	_, err := cs.categoryRepo.FindById(categoryId)
	if err != nil {
		return nil, utils.NewError("Category not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra category có con không
	hasChildren, err := cs.categoryRepo.HasChildren(categoryId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to check category children", utils.ErrCodeInternal)
	}

	if hasChildren {
		return nil, utils.NewError("Cannot delete category that has subcategories", utils.ErrCodeBadRequest)
	}

	// 3. Xóa category (soft delete vì có DeletedAt)
	if err := cs.categoryRepo.Delete(categoryId); err != nil {
		return nil, utils.WrapError(err, "Failed to delete category", utils.ErrCodeInternal)
	}

	return &dto.DeleteCategoryResponse{
		Message:    "Category deleted successfully",
		CategoryId: categoryId,
	}, nil
}
