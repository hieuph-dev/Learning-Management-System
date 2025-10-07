package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type EnrollmentRoutes struct {
	handler *handler.EnrollmentHandler
}

func NewEnrollmentRoutes(handler *handler.EnrollmentHandler) *EnrollmentRoutes {
	return &EnrollmentRoutes{
		handler: handler,
	}
}

func (er *EnrollmentRoutes) Register(r *gin.RouterGroup) {
	// Course enrollment routes
	courses := r.Group("/courses")
	{
		// Protected routes - cần authentication
		courses.Use(middleware.AuthMiddleware())
		{
			// Enroll vào course
			courses.POST("/:course_id/enroll", er.handler.EnrollCourse)

			// Kiểm tra enrollment status
			courses.GET("/course_id/:course_id/check-enrollment", er.handler.CheckEnrollment)
		}
	}

	// User enrollments routes
	enrollments := r.Group("/enrollments")
	{
		// Protected routes
		enrollments.Use(middleware.AuthMiddleware())
		{
			// Lấy danh sách enrollments của user
			enrollments.GET("/my", er.handler.GetMyEnrollments)
		}
	}
}
