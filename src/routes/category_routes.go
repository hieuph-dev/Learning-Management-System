package routes

import (
	"lms/src/handler"
	"lms/src/middleware"
	"time"

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
		categories.GET("/", middleware.CacheMiddleware(30*time.Minute), cr.handler.GetCategories)
		categories.GET("/:id", middleware.CacheMiddleware(30*time.Minute), cr.handler.GetCategoryById)
	}

	// Admin routes
	adminCategories := r.Group("/admin/categories")
	{
		adminCategories.Use(middleware.AuthMiddleware())
		adminCategories.Use(middleware.AdminMiddleware())
		{
			// Khi tạo/sửa/xóa category, xóa cache
			adminCategories.POST("/", middleware.InvalidateCachePattern("cache:/api/v1/categories*"), cr.handler.CreateCategory)
			adminCategories.PUT("/:id", middleware.InvalidateCachePattern("cache:/api/v1/categories*"), cr.handler.UpdateCategory)
			adminCategories.DELETE("/:id", middleware.InvalidateCachePattern("cache:/api/v1/categories*"), cr.handler.DeleteCategory)
		}
	}
}
