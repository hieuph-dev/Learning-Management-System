package service

import (
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
)

type adminAnalyticsService struct {
	adminAnalyticsRepo repository.AdminAnalyticsRepository
}

func NewAdminAnalyticsService(adminAnalyticsRepo repository.AdminAnalyticsRepository) AdminAnalyticsService {
	return &adminAnalyticsService{
		adminAnalyticsRepo: adminAnalyticsRepo,
	}
}

func (aas *adminAnalyticsService) GetAdminDashboard() (*dto.AdminDashboardResponse, error) {
	dashboard, err := aas.adminAnalyticsRepo.GetAdminDashboard()
	if err != nil {
		return nil, utils.WrapError(err, "failed to get admin dashboard", utils.ErrCodeInternal)
	}

	return dashboard, nil
}

func (aas *adminAnalyticsService) GetAdminRevenueAnalytics(req *dto.AdminRevenueAnalyticsRequest) (*dto.AdminRevenueAnalyticsResponse, error) {
	analytics, err := aas.adminAnalyticsRepo.GetAdminRevenueAnalytics(req)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get revenue analytics", utils.ErrCodeInternal)
	}

	return analytics, nil
}

func (aas *adminAnalyticsService) GetAdminUsersAnalytics(req *dto.AdminUsersAnalyticsRequest) (*dto.AdminUsersAnalyticsResponse, error) {
	analytics, err := aas.adminAnalyticsRepo.GetAdminUsersAnalytics(req)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get users analytics", utils.ErrCodeInternal)
	}

	return analytics, nil
}

func (aas *adminAnalyticsService) GetAdminCoursesAnalytics(req *dto.AdminCoursesAnalyticsRequest) (*dto.AdminCoursesAnalyticsResponse, error) {
	analytics, err := aas.adminAnalyticsRepo.GetAdminCoursesAnalytics(req)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get courses analytics", utils.ErrCodeInternal)
	}

	return analytics, nil
}
