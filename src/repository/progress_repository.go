package repository

import (
	"lms/src/models"

	"gorm.io/gorm"
)

type DBProgressRepository struct {
	db *gorm.DB
}

func NewDBProgressRepository(db *gorm.DB) ProgressRepository {
	return &DBProgressRepository{
		db: db,
	}
}

// CountCompletedLessons đếm số bài học đã hoàn thành của user trong course
func (pr *DBProgressRepository) CountCompletedLessons(userId, courseId uint) (int, error) {
	var count int64
	err := pr.db.Model(&models.Progress{}).
		Where("user_id = ? AND course_id = ? AND is_completed = ? AND deleted_at IS NULL", userId, courseId, true).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetCourseProgress lấy progress của tất cả lessons trong course
func (pr *DBProgressRepository) GetCourseProgress(userId, courseId uint) ([]models.Progress, error) {
	var progress []models.Progress
	err := pr.db.Where("user_id = ? AND course_id = ? AND deleted_at IS NULL", userId, courseId).
		Order("lesson_id ASC").
		Find(&progress).Error

	return progress, err
}

// UpdateProgress cập nhật hoặc tạo mới progress
func (pr *DBProgressRepository) UpdateProgress(progress *models.Progress) error {
	return pr.db.Save(progress).Error
}

// GetLessonProgress lấy progress của một lesson
func (pr *DBProgressRepository) GetLessonProgress(userId, lessonId uint) (*models.Progress, error) {
	var progress models.Progress
	err := pr.db.Where("user_id = ? AND lesson_id = ? AND deleted_at IS NULL", userId, lessonId).
		First(&progress).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &progress, nil
}
