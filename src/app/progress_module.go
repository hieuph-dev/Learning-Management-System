package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type ProgressModule struct {
	routes routes.Route
}

func NewProgressModule() *ProgressModule {
	progressRepo := repository.NewDBProgressRepository(db.DB)
	enrollmentRepo := repository.NewDBEnrollmentRepository(db.DB)
	courseRepo := repository.NewDBCourseRepository(db.DB)
	lessonRepo := repository.NewDBLessonRepository(db.DB)

	progressService := service.NewProgressService(progressRepo, enrollmentRepo, courseRepo, lessonRepo)
	progressHandler := handler.NewProgressHandler(progressService)
	progressRoutes := routes.NewProgressRoutes(progressHandler)

	return &ProgressModule{routes: progressRoutes}
}

func (pm *ProgressModule) Routes() routes.Route {
	return pm.routes
}
