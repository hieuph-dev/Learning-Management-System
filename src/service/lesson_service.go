package service

import (
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
)

type lessonService struct {
	lessonRepo repository.LessonRepository
	courseRepo repository.CourseRepository
}

func NewLessonService(lessonRepo repository.LessonRepository, courseRepo repository.CourseRepository) LessonService {
	return &lessonService{
		lessonRepo: lessonRepo,
		courseRepo: courseRepo,
	}
}

func (ls *lessonService) GetCourseLessons(userId, courseId uint) (*dto.GetCourseLessonsResponse, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := ls.courseRepo.FindById(courseId)
	if err != nil {
		// Workaround: query trực tiếp hoặc thêm method FindById
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra user đã enroll course chưa
	isEnrolled, err := ls.lessonRepo.CheckUserEnrollment(userId, courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to check enrollment", utils.ErrCodeInternal)
	}

	if !isEnrolled {
		return nil, utils.NewError("You must enroll in this course to access lessons", utils.ErrCodeForbidden)
	}

	// 3. Lấy danh sách lessons
	lessons, err := ls.lessonRepo.GetCourseLessons(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get course lessons", utils.ErrCodeInternal)
	}

	// 4. Lấy progress của user cho các lessons
	lessonIds := make([]uint, len(lessons))
	for i, lesson := range lessons {
		lessonIds[i] = lesson.Id
	}

	progressMap, err := ls.lessonRepo.GetLessonProgress(userId, lessonIds)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get lesson progress", utils.ErrCodeInternal)
	}

	// 5. Convert sang DTO
	lessonItems := make([]dto.LessonItem, len(lessons))
	for i, lesson := range lessons {
		lessonItems[i] = dto.LessonItem{
			Id:            lesson.Id,
			CourseId:      lesson.CourseId,
			Title:         lesson.Title,
			Slug:          lesson.Slug,
			Description:   lesson.Description,
			VideoURL:      lesson.VideoURL,
			VideoDuration: lesson.VideoDuration,
			LessonOrder:   lesson.LessonOrder,
			IsPreview:     lesson.IsPreview,
			IsCompleted:   progressMap[lesson.Id],
			CreatedAt:     lesson.CreatedAt,
		}
	}

	return &dto.GetCourseLessonsResponse{
		CourseId:     courseId,
		CourseTitle:  course.Title,
		Lessons:      lessonItems,
		TotalLessons: len(lessonItems),
	}, nil
}

func (ls *lessonService) GetLessonDetail(userId, courseId uint, slug string) (*dto.LessonDetail, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := ls.courseRepo.FindById(courseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra user đã enroll course chưa
	isEnrolled, err := ls.lessonRepo.CheckUserEnrollment(userId, courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to check enrollment", utils.ErrCodeInternal)
	}

	if !isEnrolled {
		return nil, utils.NewError("You must enroll in this course to access this lesson", utils.ErrCodeForbidden)
	}

	// 3. Tìm lesson theo slug và course_id
	lesson, err := ls.lessonRepo.FindLessonBySlugAndCourse(slug, courseId)
	if err != nil {
		return nil, utils.NewError("Lesson not found", utils.ErrCodeNotFound)
	}

	// 4. Lấy progress detail của lesson
	progress, err := ls.lessonRepo.GetLessonProgressDetail(userId, lesson.Id)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get lesson progress", utils.ErrCodeInternal)
	}

	// 5. Lấy previous và next lesson
	var previousLesson *dto.LessonNavigation
	var nextLesson *dto.LessonNavigation

	prevLesson, err := ls.lessonRepo.GetPreviousLesson(courseId, lesson.LessonOrder)
	if err == nil && prevLesson != nil {
		previousLesson = &dto.LessonNavigation{
			Id:    prevLesson.Id,
			Title: prevLesson.Title,
			Slug:  prevLesson.Slug,
		}
	}

	nxtLesson, err := ls.lessonRepo.GetNextLesson(courseId, lesson.LessonOrder)
	if err == nil && nxtLesson != nil {
		nextLesson = &dto.LessonNavigation{
			Id:    nextLesson.Id,
			Title: nextLesson.Title,
			Slug:  nextLesson.Slug,
		}
	}

	// 6. Convert sang DTO
	return &dto.LessonDetail{
		Id:             lesson.Id,
		CourseId:       lesson.CourseId,
		CourseTitle:    course.Title,
		Title:          lesson.Title,
		Slug:           lesson.Slug,
		Description:    lesson.Description,
		Content:        lesson.Content,
		VideoURL:       lesson.VideoURL,
		VideoDuration:  lesson.VideoDuration,
		LessonOrder:    lesson.LessonOrder,
		IsPreview:      lesson.IsPreview,
		IsCompleted:    progress.IsCompleted,
		LastPosition:   progress.LastPosition,
		WatchDuration:  progress.WatchDuration,
		CreatedAt:      lesson.CreatedAt,
		UpdatedAt:      lesson.UpdatedAt,
		PreviousLesson: previousLesson,
		NextLesson:     nextLesson,
	}, nil
}
