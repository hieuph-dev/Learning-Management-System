package service

import (
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
)

type analyticsService struct {
	analyticsRepo repository.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
	}
}

func (as *analyticsService) GetInstructorOverview(instructorId uint) (*dto.InstructorOverviewResponse, error) {
	overview, err := as.analyticsRepo.GetInstructorOverview(instructorId)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get instructor overview", utils.ErrCodeInternal)
	}

	return overview, nil
}

func (as *analyticsService) GetRevenueAnalytics(instructorId uint, req *dto.RevenueAnalyticsRequest) (*dto.RevenueAnalyticsResponse, error) {
	analytics, err := as.analyticsRepo.GetRevenueAnalytics(instructorId, req)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get revenue analytics", utils.ErrCodeInternal)
	}

	return analytics, nil
}

func (as *analyticsService) GetStudentAnalytics(instructorId uint, req *dto.StudentAnalyticsRequest) (*dto.StudentAnalyticsResponse, error) {
	analytics, err := as.analyticsRepo.GetStudentAnalytics(instructorId, req)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get student analytics", utils.ErrCodeInternal)
	}

	return analytics, nil
}
