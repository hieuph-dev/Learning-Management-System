package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule() *AuthModule {
	// Tạo repository để tương tác với database
	userRepo := repository.NewDBUserRepository(db.DB)
	passwordResetRepo := repository.NewDBPasswordResetRepository(db.DB)

	// Tạo service chứa business logic
	emailService := service.NewEmailService()
	authService := service.NewAuthService(userRepo, passwordResetRepo, emailService)

	// Tạo handler xử lý HTTP requests
	authHandler := handler.NewAuthHandler(authService)

	// Tạo routes định nghĩa các endpoint
	authRoutes := routes.NewAuthRoutes(authHandler)

	return &AuthModule{routes: authRoutes}
}

func (am *AuthModule) Routes() routes.Route {
	return am.routes
}
