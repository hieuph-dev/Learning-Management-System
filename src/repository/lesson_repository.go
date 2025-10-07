package repository

import (
	"lms/src/models"

	"gorm.io/gorm"
)

type DBLessonRepository struct {
	db *gorm.DB
}

func NewDBLessonRepository(db *gorm.DB) LessonRepository {
	return &DBLessonRepository{
		db: db,
	}
}

func (lr *DBLessonRepository) GetCourseLessons(courseId uint) ([]models.Lesson, error) {
	var lessons []models.Lesson

	err := lr.db.Where("course_id = ? AND is_published = ? AND deleted_at IS NULL", courseId, true).
		Order("lesson_order ASC").
		Find(&lessons).Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (lr *DBLessonRepository) CheckUserEnrollment(userId, courseId uint) (bool, error) {
	var count int64

	err := lr.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ? AND status = ? AND deleted_at IS NULL",
			userId, courseId, "active").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (lr *DBLessonRepository) GetLessonProgress(userId uint, lessonIds []uint) (map[uint]bool, error) {
	if len(lessonIds) == 0 {
		return make(map[uint]bool), nil
	}

	var progressList []models.Progress

	err := lr.db.Where("user_id = ? AND lesson_id IN ? AND deleted_at IS NULL",
		userId, lessonIds).
		Find(&progressList).Error

	if err != nil {
		return nil, err
	}

	progressMap := make(map[uint]bool)
	for _, progress := range progressList {
		progressMap[progress.LessonId] = progress.IsCompleted
	}

	return progressMap, nil
}

func (lr *DBLessonRepository) FindLessonBySlugAndCourse(slug string, courseId uint) (*models.Lesson, error) {
	var lesson models.Lesson

	err := lr.db.Where("slug = ? AND course_id = ? AND is_published = ? AND deleted_at IS NULL",
		slug, courseId, true).
		First(&lesson).Error

	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

func (lr *DBLessonRepository) GetLessonProgressDetail(userId, lessonId uint) (*models.Progress, error) {
	var progress models.Progress

	err := lr.db.Where("user_id = ? AND lesson_id = ? AND deleted_at IS NULL",
		userId, lessonId).
		First(&progress).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Trả về progress mặc định nếu chưa có
			return &models.Progress{
				UserId:        userId,
				LessonId:      lessonId,
				IsCompleted:   false,
				WatchDuration: 0,
				LastPosition:  0,
			}, nil
		}
		return nil, err
	}
	return &progress, nil
}

func (lr *DBLessonRepository) GetPreviousLesson(courseId uint, currentOrder int) (*models.Lesson, error) {
	var lesson models.Lesson

	err := lr.db.Where("course_id = ? AND lesson_order < ? AND is_published = ? AND deleted_at IS NULL",
		courseId, currentOrder, true).
		Order("lesson_order DESC").
		First(&lesson).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Không có lesson trước
		}
		return nil, err
	}

	return &lesson, nil
}

func (lr *DBLessonRepository) GetNextLesson(courseId uint, currentOrder int) (*models.Lesson, error) {
	var lesson models.Lesson

	err := lr.db.Where("course_id = ? AND lesson_order > ? AND is_published = ? AND deleted_at IS NULL",
		courseId, currentOrder, true).
		Order("lesson_order ASC").
		First(&lesson).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Không có lesson tiếp theo
		}
		return nil, err
	}

	return &lesson, nil
}

func (lr *DBLessonRepository) FindLessonByIds(lessonIds []uint) ([]models.Lesson, error) {
	var lessons []models.Lesson

	if len(lessonIds) == 0 {
		return lessons, nil
	}

	err := lr.db.Where("id IN ? AND deleted_at IS NULL", lessonIds).
		Find(&lessons).Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}
