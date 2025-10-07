package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type AdminModule struct {
	routes routes.Route
}

func NewAdminModule() *AdminModule {
	// Tạo repository để tương tác với database
	userRepo := repository.NewDBUserRepository(db.DB)
	courseRepo := repository.NewDBCourseRepository(db.DB)
	orderRepo := repository.NewDBOrderRepository(db.DB)
	couponRepo := repository.NewDBCouponRepository(db.DB)
	enrollmentRepo := repository.NewDBEnrollmentRepository(db.DB)
	adminAnalyticsRepo := repository.NewDBAdminAnalyticsRepository(db.DB)

	// Tạo service chứa business logic
	adminService := service.NewAdminService(userRepo, courseRepo)
	orderService := service.NewOrderService(orderRepo, courseRepo, couponRepo, enrollmentRepo)
	couponService := service.NewCouponService(couponRepo, courseRepo)
	adminAnalyticsService := service.NewAdminAnalyticsService(adminAnalyticsRepo)

	// Tạo handler xử lý HTTP requests
	adminHandler := handler.NewAdminHandler(adminService, orderService)
	couponHandler := handler.NewCouponHandler(couponService)
	adminAnalyticsHandler := handler.NewAdminAnalyticsHandler(adminAnalyticsService)

	// Tạo routes định nghĩa các endpoint
	adminRoutes := routes.NewAdminRoutes(adminHandler, couponHandler, adminAnalyticsHandler)

	return &AdminModule{routes: adminRoutes}
}

func (am *AdminModule) Routes() routes.Route {
	return am.routes
}
