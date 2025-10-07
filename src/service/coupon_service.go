package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"strings"
)

type couponService struct {
	couponRepo repository.CouponRepository
	courseRepo repository.CourseRepository
}

func NewCouponService(
	couponRepo repository.CouponRepository,
	courseRepo repository.CourseRepository,
) CouponService {
	return &couponService{
		couponRepo: couponRepo,
		courseRepo: courseRepo,
	}
}

func (cs *couponService) ValidateCoupon(req *dto.ValidateCouponRequest) (*dto.ValidateCouponResponse, error) {
	// 1. Tìm coupon
	coupon, err := cs.couponRepo.FindByCode(req.CouponCode)
	if err != nil {
		return &dto.ValidateCouponResponse{
			Valid:   false,
			Message: "Invalid coupon code",
		}, nil
	}

	// 2. Kiểm tra coupon có valid không
	if !cs.couponRepo.IsValidCoupon(coupon) {
		return &dto.ValidateCouponResponse{
			Valid:      false,
			CouponCode: req.CouponCode,
			Message:    "Coupon is expired or not available",
		}, nil
	}

	// 3. Kiểm tra course có tồn tại không
	_, err = cs.courseRepo.FindById(req.CourseId)
	if err != nil {
		return &dto.ValidateCouponResponse{
			Valid:   false,
			Message: "Course not found",
		}, nil
	}

	// 4. Kiểm tra minimum order amount
	if req.OrderTotal < coupon.MinOrderAmount {
		return &dto.ValidateCouponResponse{
			Valid:          false,
			CouponCode:     req.CouponCode,
			MinOrderAmount: coupon.MinOrderAmount,
			Message:        fmt.Sprintf("Minimum order amount for this coupon is %.2f", coupon.MinOrderAmount),
		}, nil
	}

	// 5. Tính discount amount
	discountAmount := 0.0
	if coupon.DiscountType == "percentage" {
		discountAmount = req.OrderTotal * (coupon.DiscountValue / 100)
	} else if coupon.DiscountType == "fixed" {
		discountAmount = coupon.DiscountValue
	}

	// 6. Apply max discount nếu có
	if coupon.MaxDiscountAmount != nil && discountAmount > *coupon.MaxDiscountAmount {
		discountAmount = *coupon.MaxDiscountAmount
	}

	// 7. Tính final price
	finalPrice := req.OrderTotal - discountAmount
	if finalPrice < 0 {
		finalPrice = 0
	}

	return &dto.ValidateCouponResponse{
		Valid:             true,
		CouponCode:        coupon.Code,
		DiscountType:      coupon.DiscountType,
		DiscountValue:     coupon.DiscountValue,
		DiscountAmount:    discountAmount,
		FinalPrice:        finalPrice,
		MinOrderAmount:    coupon.MinOrderAmount,
		MaxDiscountAmount: coupon.MaxDiscountAmount,
		Message:           fmt.Sprintf("Coupon applied successfully! You save %.2f", discountAmount),
	}, nil
}

func (cs *couponService) CheckCoupon(req *dto.CheckCouponRequest) (*dto.CheckCouponResponse, error) {
	// Sử dụng lại logic từ ValidateCoupon
	validateReq := &dto.ValidateCouponRequest{
		CouponCode: req.CouponCode,
		CourseId:   req.CourseId,
		OrderTotal: req.OrderTotal,
	}

	validateResp, err := cs.ValidateCoupon(validateReq)
	if err != nil {
		return nil, err
	}

	return &dto.CheckCouponResponse{
		Valid:             validateResp.Valid,
		CouponCode:        validateResp.CouponCode,
		DiscountType:      validateResp.DiscountType,
		DiscountValue:     validateResp.DiscountValue,
		DiscountAmount:    validateResp.DiscountAmount,
		FinalPrice:        validateResp.FinalPrice,
		Message:           validateResp.Message,
		MinOrderAmount:    validateResp.MinOrderAmount,
		MaxDiscountAmount: validateResp.MaxDiscountAmount,
	}, nil
}

func (cs *couponService) GetAdminCoupons(req *dto.GetAdminCouponsQueryRequest) (*dto.GetAdminCouponsResponse, error) {
	// Set defaults
	page := 1
	limit := 10
	sortBy := "desc"

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}
	if req.SortBy != "" {
		sortBy = req.SortBy
	}

	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})
	if req.IsActive != nil {
		filters["is_active"] = *req.IsActive
	}
	if req.SearchCode != "" {
		filters["search_code"] = req.SearchCode
	}

	// Get coupons
	coupons, total, err := cs.couponRepo.GetCouponsWithPagination(offset, limit, filters, "created_at", sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get coupons", utils.ErrCodeInternal)
	}

	// Convert to DTO
	couponItems := make([]dto.AdminCouponItem, len(coupons))
	for i, coupon := range coupons {
		couponItems[i] = dto.AdminCouponItem{
			Id:                coupon.Id,
			Code:              coupon.Code,
			Description:       coupon.Description,
			DiscountType:      coupon.DiscountType,
			DiscountValue:     coupon.DiscountValue,
			MinOrderAmount:    coupon.MinOrderAmount,
			MaxDiscountAmount: coupon.MaxDiscountAmount,
			UsageLimit:        coupon.UsageLimit,
			UsedCount:         coupon.UsedCount,
			ValidFrom:         coupon.ValidFrom,
			ValidTo:           coupon.ValidTo,
			IsActive:          coupon.IsActive,
			CreatedAt:         coupon.CreatedAt,
		}
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit
	pagination := dto.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return &dto.GetAdminCouponsResponse{
		Coupons:    couponItems,
		Pagination: pagination,
	}, nil
}

func (cs *couponService) CreateCoupon(req *dto.CreateCouponRequest) (*dto.CreateCouponResponse, error) {
	// Kiểm tra code đã tồn tại chưa
	existingCoupon, _ := cs.couponRepo.FindByCode(req.Code)
	if existingCoupon != nil {
		return nil, utils.NewError("Coupon code already exists", utils.ErrCodeConflict)
	}

	// Validate discount value
	if req.DiscountType == "percentage" && req.DiscountValue > 100 {
		return nil, utils.NewError("Percentage discount cannot exceed 100%", utils.ErrCodeBadRequest)
	}

	// Validate date range
	if req.ValidFrom != nil && req.ValidTo != nil {
		if req.ValidTo.Before(*req.ValidFrom) {
			return nil, utils.NewError("Valid to date must be after valid from date", utils.ErrCodeBadRequest)
		}
	}

	// Set default IsActive
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Create coupon
	coupon := &models.Coupon{
		Code:              strings.ToUpper(req.Code),
		Description:       req.Description,
		DiscountType:      req.DiscountType,
		DiscountValue:     req.DiscountValue,
		MinOrderAmount:    req.MinOrderAmount,
		MaxDiscountAmount: req.MaxDiscountAmount,
		UsageLimit:        req.UsageLimit,
		ValidFrom:         req.ValidFrom,
		ValidTo:           req.ValidTo,
		IsActive:          isActive,
		UsedCount:         0,
	}

	if err := cs.couponRepo.Create(coupon); err != nil {
		return nil, utils.WrapError(err, "Failed to create coupon", utils.ErrCodeInternal)
	}

	return &dto.CreateCouponResponse{
		Id:                coupon.Id,
		Code:              coupon.Code,
		Description:       coupon.Description,
		DiscountType:      coupon.DiscountType,
		DiscountValue:     coupon.DiscountValue,
		MinOrderAmount:    coupon.MinOrderAmount,
		MaxDiscountAmount: coupon.MaxDiscountAmount,
		UsageLimit:        coupon.UsageLimit,
		ValidFrom:         coupon.ValidFrom,
		ValidTo:           coupon.ValidTo,
		IsActive:          coupon.IsActive,
		CreatedAt:         coupon.CreatedAt,
		Message:           "Coupon created successfully",
	}, nil
}

func (cs *couponService) UpdateCoupon(couponId uint, req *dto.UpdateCouponRequest) (*dto.UpdateCouponResponse, error) {
	// Tìm coupon
	coupon, err := cs.couponRepo.FindById(couponId)
	if err != nil {
		return nil, utils.NewError("Coupon not found", utils.ErrCodeNotFound)
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.DiscountType != nil {
		updates["discount_type"] = *req.DiscountType
	}

	if req.DiscountValue != nil {
		// Validate percentage
		if req.DiscountType != nil && *req.DiscountType == "percentage" && *req.DiscountValue > 100 {
			return nil, utils.NewError("Percentage discount cannot exceed 100%", utils.ErrCodeBadRequest)
		}
		updates["discount_value"] = *req.DiscountValue
	}

	if req.MinOrderAmount != nil {
		updates["min_order_amount"] = *req.MinOrderAmount
	}

	if req.MaxDiscountAmount != nil {
		updates["max_discount_amount"] = *req.MaxDiscountAmount
	}

	if req.UsageLimit != nil {
		updates["usage_limit"] = *req.UsageLimit
	}

	if req.ValidFrom != nil {
		updates["valid_from"] = *req.ValidFrom
	}

	if req.ValidTo != nil {
		// Validate date range
		validFrom := coupon.ValidFrom
		if req.ValidFrom != nil {
			validFrom = req.ValidFrom
		}
		if validFrom != nil && req.ValidTo.Before(*validFrom) {
			return nil, utils.NewError("Valid to date must be after valid from date", utils.ErrCodeBadRequest)
		}
		updates["valid_to"] = *req.ValidTo
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Update coupon
	if err := cs.couponRepo.Update(couponId, updates); err != nil {
		return nil, utils.WrapError(err, "Failed to update coupon", utils.ErrCodeInternal)
	}

	// Get updated coupon
	updatedCoupon, err := cs.couponRepo.FindById(couponId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get updated coupon", utils.ErrCodeInternal)
	}

	return &dto.UpdateCouponResponse{
		Id:                updatedCoupon.Id,
		Code:              updatedCoupon.Code,
		Description:       updatedCoupon.Description,
		DiscountType:      updatedCoupon.DiscountType,
		DiscountValue:     updatedCoupon.DiscountValue,
		MinOrderAmount:    updatedCoupon.MinOrderAmount,
		MaxDiscountAmount: updatedCoupon.MaxDiscountAmount,
		UsageLimit:        updatedCoupon.UsageLimit,
		UsedCount:         updatedCoupon.UsedCount,
		ValidFrom:         updatedCoupon.ValidFrom,
		ValidTo:           updatedCoupon.ValidTo,
		IsActive:          updatedCoupon.IsActive,
		UpdatedAt:         updatedCoupon.UpdatedAt,
		Message:           "Coupon updated successfully",
	}, nil
}

func (cs *couponService) DeleteCoupon(couponId uint) (*dto.DeleteCouponResponse, error) {
	// Kiểm tra coupon có tồn tại không
	_, err := cs.couponRepo.FindById(couponId)
	if err != nil {
		return nil, utils.NewError("Coupon not found", utils.ErrCodeNotFound)
	}

	// Delete coupon (soft delete)
	if err := cs.couponRepo.Delete(couponId); err != nil {
		return nil, utils.WrapError(err, "Failed to delete coupon", utils.ErrCodeInternal)
	}

	return &dto.DeleteCouponResponse{
		Message: "Coupon deleted successfully",
	}, nil
}
