package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type EnrollmentModule struct {
	routes routes.Route
}

func NewEnrollmentModule() *EnrollmentModule {
	enrollmentRepo := repository.NewDBEnrollmentRepository(db.DB)
	orderRepo := repository.NewDBOrderRepository(db.DB)
	courseRepo := repository.NewDBCourseRepository(db.DB)
	couponRepo := repository.NewDBCouponRepository(db.DB)
	progressRepo := repository.NewDBProgressRepository(db.DB)

	enrollmentService := service.NewEnrollmentService(enrollmentRepo, orderRepo, courseRepo, couponRepo, progressRepo)

	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentService)

	enrollmentRoutes := routes.NewEnrollmentRoutes(enrollmentHandler)

	return &EnrollmentModule{routes: enrollmentRoutes}
}

func (em *EnrollmentModule) Routes() routes.Route {
	return em.routes
}
