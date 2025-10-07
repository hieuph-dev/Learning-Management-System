package routes

import (
	"lms/src/handler"
	"lms/src/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type CourseRoutes struct {
	handler *handler.CourseHandler
}

func NewCourseRoutes(handler *handler.CourseHandler) *CourseRoutes {
	return &CourseRoutes{
		handler: handler,
	}
}

func (cr *CourseRoutes) Register(r *gin.RouterGroup) {
	courses := r.Group("/courses")
	{
		// Public routes
		// Cache 15 phút cho danh sách courses
		courses.GET("/", middleware.CacheMiddleware(15*time.Minute), cr.handler.GetCourses)

		// Cache 10 phút cho search (vì search thay đổi thường xuyên hơn)
		courses.GET("/search", middleware.CacheMiddleware(10*time.Minute), cr.handler.SearchCourses)

		// Cache 1 giờ cho featured courses (ít thay đổi)
		courses.GET("/featured", middleware.CacheMiddleware(60*time.Minute), cr.handler.GetFeaturedCourses)

		// Cache 30 phút cho course detail
		courses.GET("/:slug", middleware.CacheMiddleware(30*time.Minute), cr.handler.GetCourseBySlug)

		// Cache 5 phút cho reviews (thay đổi khi có review mới)
		courses.GET("/course_id/:course_id/reviews", middleware.CacheMiddleware(5*time.Minute), cr.handler.GetCourseReviews)

		// Protected routes - Student can create review
		courses.Use(middleware.AuthMiddleware())
		{
			// Khi tạo review mới, xóa cache của reviews
			courses.POST("/:course_id/reviews", middleware.InvalidateCachePattern("cache:/api/v1/courses/course_id/*"), cr.handler.CreateCourseReview)
		}
	}

	// Review management routes
	reviews := r.Group("/reviews")
	{
		reviews.Use(middleware.AuthMiddleware())
		{
			// Khi update/delete review, xóa cache
			reviews.PUT("/:review_id", middleware.InvalidateCachePattern("cache:/api/v1/courses/course_id/*"), cr.handler.UpdateReview)
			reviews.DELETE("/:review_id", middleware.InvalidateCachePattern("cache:/api/v1/courses/course_id/*"), cr.handler.DeleteReview)
		}
	}
}
