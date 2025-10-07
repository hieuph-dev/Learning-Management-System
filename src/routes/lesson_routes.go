package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type LessonRoutes struct {
	handler *handler.LessonHandler
}

func NewLessonRoutes(handler *handler.LessonHandler) *LessonRoutes {
	return &LessonRoutes{
		handler: handler,
	}
}

func (lr *LessonRoutes) Register(r *gin.RouterGroup) {
	courses := r.Group("/courses")
	{
		// Protected route - cần authentication
		courses.Use(middleware.AuthMiddleware())
		{
			// Lấy danh sách lessons
			courses.GET("/course_id/:course_id/lessons", lr.handler.GetCourseLessons)

			// Lấy lesson detail
			courses.GET("/course_id/:course_id/lessons/:slug", lr.handler.GetCourseLessons)
		}
	}
}
