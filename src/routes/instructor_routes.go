package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type InstructorRoutes struct {
	handler          *handler.InstructorHandler
	analyticsHandler *handler.AnalyticsHandler
}

func NewInstructorRoutes(handler *handler.InstructorHandler, analyticsHandler *handler.AnalyticsHandler) *InstructorRoutes {
	return &InstructorRoutes{
		handler:          handler,
		analyticsHandler: analyticsHandler,
	}
}

func (ir *InstructorRoutes) Register(r *gin.RouterGroup) {
	instructor := r.Group("/instructor")
	{
		// Protected routes - cần authentication và role instructor
		instructor.Use(middleware.AuthMiddleware())
		instructor.Use(middleware.InstructorMiddleware())
		{
			// Course management
			instructor.GET("/courses", ir.handler.GetInstructorCourses)
			instructor.POST("/courses", ir.handler.CreateCourse)
			instructor.PUT("/courses/:course_id", ir.handler.UpdateCourse)
			instructor.DELETE("/courses/:course_id", ir.handler.DeleteCourse)
			instructor.GET("/courses/:course_id/students", ir.handler.GetCourseStudents)

			// Lesson management
			instructor.POST("/courses/:course_id/lessons", ir.handler.CreateLesson)
			instructor.PUT("/courses/:course_id/lessons/:id", ir.handler.UpdateLesson)
			instructor.DELETE("/courses/:course_id/lessons/:id", ir.handler.DeleteLesson)
			instructor.PUT("/lessons/:id/reorder", ir.handler.ReorderLessons)

			// Analytics endpoints
			analytics := instructor.Group("/analytics")
			{
				analytics.GET("/overview", ir.analyticsHandler.GetInstructorOverview)
				analytics.GET("/revenue", ir.analyticsHandler.GetRevenueAnalytics)
				analytics.GET("/students", ir.analyticsHandler.GetStudentAnalytics)
			}
		}
	}
}
