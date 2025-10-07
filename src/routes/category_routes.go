package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type CategoryRoutes struct {
	handler *handler.CategoryHandler
}

func NewCategoryRoutes(handler *handler.CategoryHandler) *CategoryRoutes {
	return &CategoryRoutes{
		handler: handler,
	}
}

func (cr *CategoryRoutes) Register(r *gin.RouterGroup) {
	categories := r.Group("/categories")
	{
		// Public routes
		categories.GET("/", cr.handler.GetCategories)
		categories.GET("/:id", cr.handler.GetCategoryById)
	}

	// Admin routes
	adminCategories := r.Group("/admin/categories")
	{
		adminCategories.Use(middleware.AuthMiddleware())
		adminCategories.Use(middleware.AdminMiddleware())
		{
			adminCategories.POST("/", cr.handler.CreateCategory)
			adminCategories.PUT("/:id", cr.handler.UpdateCategory)
			adminCategories.DELETE("/:id", cr.handler.DeleteCategory)
		}
	}
}
