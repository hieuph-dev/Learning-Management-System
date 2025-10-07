package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type CourseModule struct {
	routes routes.Route
}

func NewCourseModule() *CourseModule {
	courseRepo := repository.NewDBCourseRepository(db.DB)
	reviewRepo := repository.NewDBReviewRepository(db.DB)
	enrollmentRepo := repository.NewDBEnrollmentRepository(db.DB)

	courseService := service.NewCourseService(courseRepo)
	reviewService := service.NewReviewService(reviewRepo, courseRepo, enrollmentRepo)

	courseHandler := handler.NewCourseHandler(courseService, reviewService)

	courseRoutes := routes.NewCourseRoutes(courseHandler)

	return &CourseModule{routes: courseRoutes}
}

func (cm *CourseModule) Routes() routes.Route {
	return cm.routes
}
