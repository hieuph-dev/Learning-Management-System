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

type CouponHandler struct {
	couponService service.CouponService
}

func NewCouponHandler(couponService service.CouponService) *CouponHandler {
	return &CouponHandler{
		couponService: couponService,
	}
}

// POST /api/v1/coupons/check - Check coupon (Public)
func (ch *CouponHandler) CheckCoupon(ctx *gin.Context) {
	var req dto.CheckCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	response, err := ch.couponService.CheckCoupon(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/admin/coupons - Get all coupons (Admin)
func (ch *CouponHandler) GetAdminCoupons(ctx *gin.Context) {
	var req dto.GetAdminCouponsQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	response, err := ch.couponService.GetAdminCoupons(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// POST /api/v1/admin/coupons - Create coupon (Admin)
func (ch *CouponHandler) CreateCoupon(ctx *gin.Context) {
	var req dto.CreateCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	response, err := ch.couponService.CreateCoupon(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusCreated, response)
}

// PUT /api/v1/admin/coupons/:id - Update coupon (Admin)
func (ch *CouponHandler) UpdateCoupon(ctx *gin.Context) {
	couponIdParam := ctx.Param("id")
	if couponIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Coupon Id is required", utils.ErrCodeBadRequest))
		return
	}

	couponId, err := strconv.ParseUint(couponIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid coupon Id format", utils.ErrCodeBadRequest))
		return
	}

	var req dto.UpdateCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	response, err := ch.couponService.UpdateCoupon(uint(couponId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// DELETE /api/v1/admin/coupons/:id - Delete coupon (Admin)
func (ch *CouponHandler) DeleteCoupon(ctx *gin.Context) {
	couponIdParam := ctx.Param("id")
	if couponIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Coupon Id is required", utils.ErrCodeBadRequest))
		return
	}

	couponId, err := strconv.ParseUint(couponIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid coupon Id format", utils.ErrCodeBadRequest))
		return
	}

	response, err := ch.couponService.DeleteCoupon(uint(couponId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// Thêm vào cuối file coupon_handler.go

// Implement Route interface
func (ch *CouponHandler) Register(r *gin.RouterGroup) {
	// Public coupon routes
	coupons := r.Group("/coupons")
	{
		coupons.POST("/check", ch.CheckCoupon)
	}
}
