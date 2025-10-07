package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type CategoryModule struct {
	routes routes.Route
}

func NewCategoryModule() *CategoryModule {
	categoryRepo := repository.NewDBCategoryRepository(db.DB)

	categoryService := service.NewCategoryService(categoryRepo)

	categoryHandler := handler.NewCategoryHandler(categoryService)

	categoryRoutes := routes.NewCategoryRoutes(categoryHandler)

	return &CategoryModule{routes: categoryRoutes}
}

func (cm *CategoryModule) Routes() routes.Route {
	return cm.routes
}
