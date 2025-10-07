package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewDBUserRepository(db.DB)

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{routes: userRoutes}
}

func (um *UserModule) Routes() routes.Route {
	return um.routes
}
