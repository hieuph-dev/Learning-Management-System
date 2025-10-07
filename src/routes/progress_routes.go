package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type ProgressRoutes struct {
	handler *handler.ProgressHandler
}

func NewProgressRoutes(handler *handler.ProgressHandler) *ProgressRoutes {
	return &ProgressRoutes{
		handler: handler,
	}
}

func (pr *ProgressRoutes) Register(r *gin.RouterGroup) {
	// Enrollment progress routes
	enrollments := r.Group("/enrollments")
	{
		// Protected routes - cần authentication
		enrollments.Use(middleware.AuthMiddleware())
		{
			// Lấy learning progress của course
			enrollments.GET("/:course_id/progress", pr.handler.GetCourseProgress)
		}
	}

	// Lesson progress routes
	progress := r.Group("/progress")
	{
		progress.Use(middleware.AuthMiddleware())
		{
			// Đánh dấu lesson hoàn thành
			progress.POST("/:lesson_id/complete", pr.handler.CompleteLesson)

			// Cập nhật vị trí video
			progress.PUT("/:lesson_id/position", pr.handler.UpdateLessonPosition)
		}
	}
}
