package handler

import (
	"lms/src/dto"
	"lms/src/service"
	"lms/src/utils"
	"lms/src/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminAnalyticsHandler struct {
	service service.AdminAnalyticsService
}

func NewAdminAnalyticsHandler(service service.AdminAnalyticsService) *AdminAnalyticsHandler {
	return &AdminAnalyticsHandler{
		service: service,
	}
}

// GET /api/v1/admin/analytics/dashboard
func (aah *AdminAnalyticsHandler) GetAdminDashboard(ctx *gin.Context) {
	// Gọi service
	response, err := aah.service.GetAdminDashboard()
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/admin/analytics/revenue
func (aah *AdminAnalyticsHandler) GetAdminRevenueAnalytics(ctx *gin.Context) {
	// Parse query parameters
	var req dto.AdminRevenueAnalyticsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service
	response, err := aah.service.GetAdminRevenueAnalytics(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/admin/analytics/users
func (aah *AdminAnalyticsHandler) GetAdminUsersAnalytics(ctx *gin.Context) {
	// Parse query parameters
	var req dto.AdminUsersAnalyticsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service
	response, err := aah.service.GetAdminUsersAnalytics(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/admin/analytics/courses
func (aah *AdminAnalyticsHandler) GetAdminCoursesAnalytics(ctx *gin.Context) {
	// Parse query parameters
	var req dto.AdminCoursesAnalyticsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service
	response, err := aah.service.GetAdminCoursesAnalytics(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
