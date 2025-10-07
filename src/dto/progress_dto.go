package dto

import "time"

// GetCourseProgressResponse - Response chi tiết progress của course
type GetCourseProgressResponse struct {
	CourseId           uint                 `json:"course_id"`
	CourseTitle        string               `json:"course_title"`
	IsEnrolled         bool                 `json:"is_enrolled"`
	EnrolledAt         *time.Time           `json:"enrolled_at,omitempty"`
	ProgressPercentage float64              `json:"progress_percentage"`
	TotalLessons       int                  `json:"total_lessons"`
	CompletedLessons   int                  `json:"completed_lessons"`
	TotalDuration      int                  `json:"total_duration"`   // Tổng thời lượng (giây)
	WatchedDuration    int                  `json:"watched_duration"` // Đã xem (giây)
	LastAccessedAt     *time.Time           `json:"last_accessed_at,omitempty"`
	Status             string               `json:"status"` // active, completed, dropped
	Lessons            []LessonProgressItem `json:"lessons"`
}

// LessonProgressItem - Progress của từng lesson
type LessonProgressItem struct {
	LessonId        uint       `json:"lesson_id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	LessonOrder     int        `json:"lesson_order"`
	VideoDuration   int        `json:"video_duration"` // Tổng thời lượng video (giây)
	IsCompleted     bool       `json:"is_completed"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	WatchDuration   int        `json:"watch_duration"`   // Đã xem (giây)
	LastPosition    int        `json:"last_position"`    // Vị trí cuối cùng (giây)
	ProgressPercent float64    `json:"progress_percent"` // % hoàn thành lesson này
}

// Request để đánh dấu lesson hoàn thành
type CompleteLessonRequest struct {
	WatchDuration int `json:"watch_duration" binding:"required,min=0"`
}

type CompleteLessonResponse struct {
	LessonId      uint      `json:"lesson_id"`
	CourseId      uint      `json:"course_id"`
	IsCompleted   bool      `json:"is_completed"`
	CompletedAt   time.Time `json:"completed_at"`
	WatchDuration int       `json:"watch_duration"`
	Message       string    `json:"message"`
}

// Request để cập nhật vị trí video
type UpdateLessonPositionRequest struct {
	LastPosition  int `json:"last_position" binding:"required,min=0"`
	WatchDuration int `json:"watch_duration" binding:"required,min=0"`
}

type UpdateLessonPositionResponse struct {
	LessonId      uint   `json:"lesson_id"`
	LastPosition  int    `json:"last_position"`
	WatchDuration int    `json:"watch_duration"`
	Message       string `json:"message"`
}
