package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"time"
)

type progressService struct {
	progressRepo   repository.ProgressRepository
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	lessonRepo     repository.LessonRepository
}

func NewProgressService(
	progressRepo repository.ProgressRepository,
	enrollmentRepo repository.EnrollmentRepository,
	courseRepo repository.CourseRepository,
	lessonRepo repository.LessonRepository,
) ProgressService {
	return &progressService{
		progressRepo:   progressRepo,
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		lessonRepo:     lessonRepo,
	}
}

func (ps *progressService) GetCourseProgress(userId, courseId uint) (*dto.GetCourseProgressResponse, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := ps.courseRepo.FindById(courseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra user đã enroll chưa
	enrollment, isEnrolled := ps.enrollmentRepo.CheckEnrollment(userId, courseId)
	if !isEnrolled {
		return nil, utils.NewError("You are not enrolled in this course", utils.ErrCodeForbidden)
	}

	// 3. Lấy danh sách lessons của course
	lessons, err := ps.lessonRepo.GetCourseLessons(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get course lessons", utils.ErrCodeInternal)
	}

	// 4. Lấy progress của tất cả lessons
	progressMap := make(map[uint]*dto.LessonProgressItem)
	courseProgress, err := ps.progressRepo.GetCourseProgress(userId, courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get progress", utils.ErrCodeInternal)
	}

	// Map progress theo lesson_id
	progressDataMap := make(map[uint]struct {
		isCompleted   bool
		completedAt   *time.Time
		watchDuration int
		lastPosition  int
	})

	for _, p := range courseProgress {
		progressDataMap[p.LessonId] = struct {
			isCompleted   bool
			completedAt   *time.Time
			watchDuration int
			lastPosition  int
		}{
			isCompleted:   p.IsCompleted,
			completedAt:   p.CompletedAt,
			watchDuration: p.WatchDuration,
			lastPosition:  p.LastPosition,
		}
	}

	// 5. Tính toán progress cho từng lesson
	totalDuration := 0   // tổng thời lượng video của tất cả lessons.
	watchedDuration := 0 // tổng thời lượng mà user đã xem.
	completedCount := 0  // số bài học user đã hoàn thành.

	lessonItems := make([]dto.LessonProgressItem, 0, len(lessons))

	for _, lesson := range lessons {
		totalDuration += lesson.VideoDuration

		progressPercent := 0.0
		isCompleted := false
		var completedAt *time.Time = nil
		watchDuration := 0
		lastPosition := 0

		// Kiểm tra có progress không
		if p, exists := progressDataMap[lesson.Id]; exists {
			isCompleted = p.isCompleted
			completedAt = p.completedAt
			watchDuration = p.watchDuration
			lastPosition = p.lastPosition
			watchedDuration += watchDuration

			if isCompleted {
				completedCount++
				progressPercent = 100.0
			} else if lesson.VideoDuration > 0 {
				progressPercent = float64(watchDuration) / float64(lesson.VideoDuration) * 100
				if progressPercent > 100 {
					progressPercent = 100
				}
			}
		}

		lessonItems = append(lessonItems, dto.LessonProgressItem{
			LessonId:        lesson.Id,
			Title:           lesson.Title,
			Slug:            lesson.Slug,
			LessonOrder:     lesson.LessonOrder,
			VideoDuration:   lesson.VideoDuration,
			IsCompleted:     isCompleted,
			CompletedAt:     completedAt,
			WatchDuration:   watchDuration,
			LastPosition:    lastPosition,
			ProgressPercent: progressPercent,
		})

		progressMap[lesson.Id] = &lessonItems[len(lessonItems)-1]
	}

	// 6. Tính progress percentage tổng thể
	overallProgress := 0.0
	if len(lessons) > 0 {
		overallProgress = float64(completedCount) / float64(len(lessons)) * 100
	}

	// 7. Trả về response
	return &dto.GetCourseProgressResponse{
		CourseId:           course.Id,
		CourseTitle:        course.Title,
		IsEnrolled:         true,
		EnrolledAt:         &enrollment.EnrolledAt,
		ProgressPercentage: overallProgress,
		TotalLessons:       len(lessons),
		CompletedLessons:   completedCount,
		TotalDuration:      totalDuration,
		WatchedDuration:    watchedDuration,
		LastAccessedAt:     enrollment.LastAccessedAt,
		Status:             enrollment.Status,
		Lessons:            lessonItems,
	}, nil
}

func (ps *progressService) CompleteLesson(userId, lessonId uint, req *dto.CompleteLessonRequest) (*dto.CompleteLessonResponse, error) {
	// 1. Lấy thông tin lesson
	lessons, err := ps.lessonRepo.FindLessonByIds([]uint{lessonId})
	if err != nil || len(lessons) == 0 {
		return nil, utils.NewError("Lesson not found", utils.ErrCodeNotFound)
	}
	lesson := lessons[0]

	// 2. Kiểm tra user đã enroll course chưa
	_, isEnrolled := ps.enrollmentRepo.CheckEnrollment(userId, lesson.CourseId)
	if !isEnrolled {
		return nil, utils.NewError("You are not enrolled in this course", utils.ErrCodeForbidden)
	}

	// 3. Lấy hoặc tạo progress record
	progress, err := ps.progressRepo.GetLessonProgress(userId, lessonId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get lesson progress", utils.ErrCodeInternal)
	}

	// Nếu chưa có progress, tạo mới
	if progress == nil {
		now := time.Now()
		progress = &models.Progress{
			UserId:        userId,
			LessonId:      lessonId,
			CourseId:      lesson.CourseId,
			IsCompleted:   true,
			CompletedAt:   &now,
			WatchDuration: req.WatchDuration,
			LastPosition:  lesson.VideoDuration, // Set to end
		}
	} else {
		// Cập nhật progress hiện tại
		now := time.Now()
		progress.IsCompleted = true
		progress.CompletedAt = &now
		progress.WatchDuration = req.WatchDuration
		progress.LastPosition = lesson.VideoDuration
	}

	// 4. Lưu progress
	if err := ps.progressRepo.UpdateProgress(progress); err != nil {
		return nil, utils.WrapError(err, "Failed to update progress", utils.ErrCodeInternal)
	}

	// 5. Cập nhật enrollment progress percentage
	if err := ps.updateEnrollmentProgress(userId, lesson.CourseId); err != nil {
		// Log error nhưng không fail request
		fmt.Printf("Failed to update enrollment progress: %v\n", err)
	}

	return &dto.CompleteLessonResponse{
		LessonId:      lessonId,
		CourseId:      lesson.CourseId,
		IsCompleted:   true,
		CompletedAt:   *progress.CompletedAt,
		WatchDuration: progress.WatchDuration,
		Message:       "Lesson completed successfully",
	}, nil
}

func (ps *progressService) UpdateLessonPosition(userId, lessonId uint, req *dto.UpdateLessonPositionRequest) (*dto.UpdateLessonPositionResponse, error) {
	// 1. Lấy thông tin lesson
	lessons, err := ps.lessonRepo.FindLessonByIds([]uint{lessonId})
	if err != nil || len(lessons) == 0 {
		return nil, utils.NewError("Lesson not found", utils.ErrCodeNotFound)
	}
	lesson := lessons[0]

	// 2. Kiểm tra user đã enroll course chưa
	_, isEnrolled := ps.enrollmentRepo.CheckEnrollment(userId, lesson.CourseId)
	if !isEnrolled {
		return nil, utils.NewError("You are not enrolled in this course", utils.ErrCodeForbidden)
	}

	// 3. Validate position không vượt quá video duration
	if req.LastPosition > lesson.VideoDuration {
		req.LastPosition = lesson.VideoDuration
	}

	// 4. Lấy hoặc tạo progress record
	progress, err := ps.progressRepo.GetLessonProgress(userId, lessonId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get lesson progress", utils.ErrCodeInternal)
	}

	if progress == nil {
		// Tạo mới progress
		progress = &models.Progress{
			UserId:        userId,
			LessonId:      lessonId,
			CourseId:      lesson.CourseId,
			IsCompleted:   false,
			WatchDuration: req.WatchDuration,
			LastPosition:  req.LastPosition,
		}
	} else {
		// Cập nhật progress hiện tại
		progress.WatchDuration = req.WatchDuration
		progress.LastPosition = req.LastPosition
	}

	// 5. Lưu progress
	if err := ps.progressRepo.UpdateProgress(progress); err != nil {
		return nil, utils.WrapError(err, "Failed to update progress", utils.ErrCodeInternal)
	}

	return &dto.UpdateLessonPositionResponse{
		LessonId:      lessonId,
		LastPosition:  progress.LastPosition,
		WatchDuration: progress.WatchDuration,
		Message:       "Video position updated successfully",
	}, nil
}

// Tiếp theo hàm updateEnrollmentProgress
func (ps *progressService) updateEnrollmentProgress(userId, courseId uint) error {
	// Đếm số lessons đã hoàn thành
	completedCount, err := ps.progressRepo.CountCompletedLessons(userId, courseId)
	if err != nil {
		return err
	}

	// Lấy tổng số lessons
	lessons, err := ps.lessonRepo.GetCourseLessons(courseId)
	if err != nil {
		return err
	}

	totalLessons := len(lessons)
	if totalLessons == 0 {
		return nil
	}

	// Tính progress percentage
	progressPercentage := float64(completedCount) / float64(totalLessons) * 100

	// Cập nhật enrollment
	enrollment, exists := ps.enrollmentRepo.CheckEnrollment(userId, courseId)
	if !exists {
		return fmt.Errorf("enrollment not found")
	}

	// Update enrollment progress
	updates := map[string]interface{}{
		"progress_percentage": progressPercentage,
		"last_accessed_at":    time.Now(),
	}

	// Nếu hoàn thành 100%, cập nhật status
	if progressPercentage >= 100 {
		updates["status"] = "completed"
		now := time.Now()
		updates["completed_at"] = now
	}

	return ps.enrollmentRepo.UpdateEnrollmentProgress(enrollment.Id, updates)
}
