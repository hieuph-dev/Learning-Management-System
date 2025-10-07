package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRoutes(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}

func (ur *UserRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		// Protected routes - cần authentication
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", ur.handler.GetProfile)
			users.PUT("/profile", ur.handler.UpdateProfile)
			users.PUT("/change-password", ur.handler.ChangePassword)
			users.POST("/upload-avatar", ur.handler.UploadAvatar)
		}
	}
}
