package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type LessonModule struct {
	routes routes.Route
}

func NewLessonModule() *LessonModule {
	lessonRepo := repository.NewDBLessonRepository(db.DB)
	courseRepo := repository.NewDBCourseRepository(db.DB)

	lessonService := service.NewLessonService(lessonRepo, courseRepo)

	lessonHandler := handler.NewLessonHandler(lessonService)

	lessonRoutes := routes.NewLessonRoutes(lessonHandler)

	return &LessonModule{routes: lessonRoutes}
}

func (lm *LessonModule) Routes() routes.Route {
	return lm.routes
}
