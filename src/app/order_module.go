package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type OrderModule struct {
	routes routes.Route
}

func NewOrderModule() *OrderModule {
	orderRepo := repository.NewDBOrderRepository(db.DB)
	courseRepo := repository.NewDBCourseRepository(db.DB)
	couponRepo := repository.NewDBCouponRepository(db.DB)
	enrollmentRepo := repository.NewDBEnrollmentRepository(db.DB)

	orderService := service.NewOrderService(orderRepo, courseRepo, couponRepo, enrollmentRepo)
	couponService := service.NewCouponService(couponRepo, courseRepo)

	orderHandler := handler.NewOrderHandler(orderService, couponService)

	orderRoutes := routes.NewOrderRoutes(orderHandler)

	return &OrderModule{routes: orderRoutes}
}

func (om *OrderModule) Routes() routes.Route {
	return om.routes
}
