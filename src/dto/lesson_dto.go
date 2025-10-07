package dto

import "time"

type LessonItem struct {
	Id            uint      `json:"id"`
	CourseId      uint      `json:"course_id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	VideoURL      string    `json:"video_url"`
	VideoDuration int       `json:"video_duration"`
	LessonOrder   int       `json:"lesson_order"`
	IsPreview     bool      `json:"is_preview"`
	IsCompleted   bool      `json:"is_completed"` // Trạng thái hoàn thành của student
	CreatedAt     time.Time `json:"created_at"`
}

type GetCourseLessonsResponse struct {
	CourseId     uint         `json:"course_id"`
	CourseTitle  string       `json:"course_title"`
	Lessons      []LessonItem `json:"lessons"`
	TotalLessons int          `json:"total_lessons"`
}

// DTO mới cho lesson detail
type LessonDetail struct {
	Id            uint      `json:"id"`
	CourseId      uint      `json:"course_id"`
	CourseTitle   string    `json:"course_title"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	Content       string    `json:"content"`
	VideoURL      string    `json:"video_url"`
	VideoDuration int       `json:"video_duration"`
	LessonOrder   int       `json:"lesson_order"`
	IsPreview     bool      `json:"is_preview"`
	IsCompleted   bool      `json:"is_completed"`
	LastPosition  int       `json:"last_position"`
	WatchDuration int       `json:"watch_duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Navigation
	PreviousLesson *LessonNavigation `json:"previous_lesson,omitempty"`
	NextLesson     *LessonNavigation `json:"next_lesson,omitempty"`
}

type LessonNavigation struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}
