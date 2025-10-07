package repository

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"strings"

	"gorm.io/gorm"
)

type DBInstructorRepository struct {
	db *gorm.DB
}

func NewDBInstructorRepository(db *gorm.DB) InstructorRepository {
	return &DBInstructorRepository{
		db: db,
	}
}

func (ir *DBInstructorRepository) GetInstructorCourses(instructorId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Course, int, error) {
	var courses []models.Course
	var total int64

	query := ir.db.Model(&models.Course{}).
		Preload("Category").
		Where("instructor_id = ?", instructorId)

	// Apply filters
	for field, value := range filters {
		if field == "search" {
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", orderBy, strings.ToUpper(sortBy)))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, int(total), nil
}

func (ir *DBInstructorRepository) CreateCourse(course *models.Course) error {
	return ir.db.Create(course).Error
}

func (ir *DBInstructorRepository) FindCourseBySlug(slug string) (*models.Course, bool) {
	var course models.Course
	if err := ir.db.Where("slug = ?", slug).First(&course).Error; err != nil {
		return nil, false
	}
	return &course, true
}

func (ir *DBInstructorRepository) FindCourseById(courseId uint) (*models.Course, error) {
	var course models.Course
	if err := ir.db.Preload("Category").Where("id = ?", courseId).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (ir *DBInstructorRepository) FindCourseByIdAndInstructor(courseId, instructorId uint) (*models.Course, error) {
	var course models.Course
	if err := ir.db.Preload("Category").
		Where("id = ? AND instructor_id = ?", courseId, instructorId).
		First(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (ir *DBInstructorRepository) UpdateCourse(courseId uint, updates map[string]interface{}) error {
	return ir.db.Model(&models.Course{}).Where("id = ?", courseId).Updates(updates).Error
}

func (ir *DBInstructorRepository) DeleteCourse(courseId uint) error {
	return ir.db.Delete(&models.Course{}, courseId).Error
}

func (ir *DBInstructorRepository) CountEnrollmentsByCourse(courseId uint) (int64, error) {
	var count int64
	if err := ir.db.Model(&models.Enrollment{}).
		Where("course_id = ?", courseId).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (ir *DBInstructorRepository) GetCourseStudents(courseId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Enrollment, int, error) {
	var enrollments []models.Enrollment
	var total int64

	query := ir.db.Model(&models.Enrollment{}).
		Preload("User").
		Where("course_id", courseId)

	// Apply filters
	for field, value := range filters {
		if field == "search" {
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Joins("JOIN users ON users.id = enrollments.user_id").
				Where("users.username ILIKE ? OR users.email ILIKE ? OR users.full_name ILIKE ?",
					searchTerm, searchTerm, searchTerm)
		} else {
			query = query.Where(fmt.Sprintf("enrollments.%s = ?", field), value)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("enrollments.%s %s", orderBy, strings.ToUpper(sortBy)))
	} else {
		query = query.Order("enrollments.enrolled_at DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, int(total), nil
}

func (ir *DBInstructorRepository) GetStudentStatistics(courseId uint) (*dto.StudentStatistics, error) {
	var stats dto.StudentStatistics

	// Total students
	var totalStudents int64
	if err := ir.db.Model(&models.Enrollment{}).
		Where("course_id = ?", courseId).
		Count(&totalStudents).Error; err != nil {
		return nil, err
	}
	stats.TotalStudents = int(totalStudents)

	// Active students
	var activeStudents int64
	if err := ir.db.Model(&models.Enrollment{}).
		Where("course_id = ? AND status = ?", courseId, "active").
		Count(&activeStudents).Error; err != nil {
		return nil, err
	}
	stats.ActiveStudents = int(activeStudents)

	// Completed students
	var completedStudents int64
	if err := ir.db.Model(&models.Enrollment{}).
		Where("course_id = ? AND status = ?", courseId, "completed").
		Count(&completedStudents).Error; err != nil {
		return nil, err
	}
	stats.CompletedStudents = int(completedStudents)

	// Dropped students
	var droppedStudents int64
	if err := ir.db.Model(&models.Enrollment{}).
		Where("course_id = ? AND status = ?", courseId, "dropped").
		Count(&droppedStudents).Error; err != nil {
		return nil, err
	}
	stats.DroppedStudents = int(droppedStudents)

	// Average progress
	var avgProgress struct {
		Avg float64
	}
	if err := ir.db.Model(&models.Enrollment{}).
		Select("COALESCE(AVG(progress_percentage), 0) as avg").
		Where("course_id = ?", courseId).
		Scan(&avgProgress).Error; err != nil {
		return nil, err
	}
	stats.AverageProgress = avgProgress.Avg

	return &stats, nil
}

func (ir *DBInstructorRepository) CreateLesson(lesson *models.Lesson) error {
	return ir.db.Create(lesson).Error
}

func (ir *DBInstructorRepository) FindLessonBySlug(slug string, courseId uint) (*models.Lesson, bool) {
	var lesson models.Lesson
	err := ir.db.Where("slug = ? AND course_id = ? AND deleted_at IS NULL", slug, courseId).
		First(&lesson).Error

	if err != nil {
		return nil, false
	}
	return &lesson, true
}

func (ir *DBInstructorRepository) CheckLessonOrderExists(courseId uint, lessonOrder int) (bool, error) {
	var count int64
	err := ir.db.Model(&models.Lesson{}).
		Where("course_id = ? AND lesson_order = ? AND deleted_at IS NULL", courseId, lessonOrder).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ir *DBInstructorRepository) FindLessonByIdAndCourse(lessonId, courseId uint) (*models.Lesson, error) {
	var lesson models.Lesson
	err := ir.db.Where("id = ? AND course_id = ? AND deleted_at IS NULL", lessonId, courseId).
		First(&lesson).Error

	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (ir *DBInstructorRepository) UpdateLesson(lessonId uint, updates map[string]interface{}) error {
	return ir.db.Model(&models.Lesson{}).
		Where("id = ?", lessonId).
		Updates(updates).Error
}

func (ir *DBInstructorRepository) DeleteLesson(lessonId uint) error {
	return ir.db.Where("id = ?", lessonId).
		Delete(&models.Lesson{}).Error
}

func (ir *DBInstructorRepository) CheckLessonOrderExistsExcept(courseId uint, lessonOrder int, excludeId uint) (bool, error) {
	var count int64
	err := ir.db.Model(&models.Lesson{}).
		Where("course_id = ? AND lesson_order = ? AND id != ? AND deleted_at IS NULL",
			courseId, lessonOrder, excludeId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ir *DBInstructorRepository) FindLessonsByIds(lessonIds []uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := ir.db.Where("id IN ? AND deleted_at IS NULL", lessonIds).
		Find(&lessons).Error

	if err != nil {
		return nil, err
	}
	return lessons, nil
}

func (ir *DBInstructorRepository) UpdateLessonOrder(lessonId uint, newOrder int) error {
	return ir.db.Model(&models.Lesson{}).
		Where("id = ?", lessonId).
		Update("lesson_order", newOrder).Error
}

func (ir *DBInstructorRepository) BeginTransaction() *gorm.DB {
	return ir.db.Begin()
}
