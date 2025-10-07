package repository

import (
	"fmt"
	"lms/src/models"

	"gorm.io/gorm"
)

type DBEnrollmentRepository struct {
	db *gorm.DB
}

func NewDBEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &DBEnrollmentRepository{
		db: db,
	}
}

func (er *DBEnrollmentRepository) Create(enrollment *models.Enrollment) error {
	return er.db.Create(enrollment).Error
}

func (er *DBEnrollmentRepository) CheckEnrollment(userId, courseId uint) (*models.Enrollment, bool) {
	var enrollment models.Enrollment
	err := er.db.Where("user_id = ? AND course_id = ? AND deleted_at IS NULL", userId, courseId).
		First(&enrollment).Error

	if err != nil {
		return nil, false
	}

	return &enrollment, true
}

func (er *DBEnrollmentRepository) GetUserEnrollments(userId uint, offset, limit int, filters map[string]interface{}) ([]models.Enrollment, int, error) {
	var enrollments []models.Enrollment
	var total int64

	query := er.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND deleted_at IS NULL", userId)

	// Apply filters
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get enrollments with Course preloaded
	// ✅ Thêm Preload để load course và instructor info
	if err := query.
		Preload("Course").            // Load course info
		Preload("Course.Instructor"). // Load instructor info
		Order("enrolled_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, int(total), nil
}

func (er *DBEnrollmentRepository) CompleteEnrollment(enrollmentId uint) error {
	updates := map[string]interface{}{
		"status":              "completed",
		"progress_percentage": 100,
		"completed_at":        gorm.Expr("NOW()"),
	}

	return er.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentId).
		Updates(updates).Error
}

func (er *DBEnrollmentRepository) UpdateEnrollmentProgress(enrollmentId uint, updates map[string]interface{}) error {
	return er.db.Model(&models.Enrollment{}).
		Where("id = ?", enrollmentId).
		Updates(updates).Error
}

func (er *DBEnrollmentRepository) CheckUserEnrollment(userId, courseId uint) (bool, error) {
	var count int64
	err := er.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ? AND status = ? AND deleted_at IS NULL",
			userId, courseId, "active").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
