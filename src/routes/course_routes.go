package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

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
		courses.GET("/", cr.handler.GetCourses)
		courses.GET("/search", cr.handler.SearchCourses)
		courses.GET("/featured", cr.handler.GetFeaturedCourses)
		courses.GET("/:slug", cr.handler.GetCourseBySlug)
		courses.GET("/course_id/:course_id/reviews", cr.handler.GetCourseReviews)

		// Protected routes - Student can create review
		courses.Use(middleware.AuthMiddleware())
		{
			courses.POST("/:course_id/reviews", cr.handler.CreateCourseReview)
		}
	}

	// Review management routes
	reviews := r.Group("/reviews")
	{
		reviews.Use(middleware.AuthMiddleware())
		{
			reviews.PUT("/:review_id", cr.handler.UpdateReview)
			reviews.DELETE("/:review_id", cr.handler.DeleteReview)
		}
	}
}
