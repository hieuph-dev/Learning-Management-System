package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type InstructorModule struct {
	routes routes.Route
}

func NewInstructorModule() *InstructorModule {
	instructorRepo := repository.NewDBInstructorRepository(db.DB)
	categoryRepo := repository.NewDBCategoryRepository(db.DB)
	analyticsRepo := repository.NewDBAnalyticsRepository(db.DB)

	instructorService := service.NewInstructorService(instructorRepo, categoryRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)

	instructorHandler := handler.NewInstructorHandler(instructorService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	instructorRoutes := routes.NewInstructorRoutes(instructorHandler, analyticsHandler)

	return &InstructorModule{routes: instructorRoutes}
}

func (im *InstructorModule) Routes() routes.Route {
	return im.routes
}
