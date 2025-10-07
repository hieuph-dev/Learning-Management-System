package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	handler *handler.AuthHandler
}

func NewAuthRoutes(handler *handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		handler: handler,
	}
}

// Dang ki route auth
func (ar *AuthRoutes) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		// Public routes - không cần authentication
		auth.POST("/register", ar.handler.Register)
		auth.POST("/login", ar.handler.Login)
		auth.POST("/refresh", ar.handler.RefreshToken)
		auth.POST("/forgot-password", ar.handler.ForgotPassword)
		auth.POST("/reset-password", ar.handler.ResetPassword)

		// Protected routes - cần authentication
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		auth.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", ar.handler.GetProfile)
			protected.POST("/logout", ar.handler.Logout)
		}
	}
}
