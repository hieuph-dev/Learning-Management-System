package repository

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"strings"

	"gorm.io/gorm"
)

type DBCourseRepository struct {
	db *gorm.DB
}

func NewDBCourseRepository(db *gorm.DB) CourseRepository {
	return &DBCourseRepository{
		db: db,
	}
}

func (cr *DBCourseRepository) GetCoursesWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Course, int, error) {
	var courses []models.Course
	var total int64

	query := cr.db.Model(&models.Course{}).
		Preload("Instructor").
		Preload("Category").
		Where("deleted_at IS NULL")

	// Apply filters
	for field, value := range filters {
		switch field {
		case "search":
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Where("title ILIKE ? OR description ILIKE ? OR short_desc ILIKE ?", searchTerm, searchTerm, searchTerm)
		case "min_price":
			query = query.Where("price >= ?", value)
		case "max_price":
			query = query.Where("price <= ?", value)
		case "price_range":
			// Handle price range if needed
		default:
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

func (cr *DBCourseRepository) SearchCourses(query string, offset, limit int, filters map[string]interface{}, sortBy, order string) ([]models.Course, int, error) {
	var courses []models.Course
	var total int64

	dbQuery := cr.db.Model(&models.Course{}).
		Preload("Instructor").
		Preload("Category").
		Where("deleted_at IS NULL AND status = ?", "published")

	// Full text search
	searchTerm := fmt.Sprintf("%%%s%%", query)
	dbQuery = dbQuery.Where(
		"title ILIKE ? OR description ILIKE ? OR short_desc ILIKE ? OR requirements ILIKE ? OR what_you_learn ILIKE ?",
		searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
	)

	// Apply filters
	for field, value := range filters {
		switch field {
		case "min_price":
			dbQuery = dbQuery.Where("price >= ?", value)
		case "max_price":
			dbQuery = dbQuery.Where("price <= ?", value)
		default:
			dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	// Count total records
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderClause := "created_at DESC" // default
	switch sortBy {
	case "relevance":
		// For now, use created_at as relevance. In production, you might use full-text search ranking
		orderClause = "created_at DESC"
	case "price":
		orderClause = fmt.Sprintf("price %s", strings.ToUpper(order))
	case "rating_avg":
		orderClause = fmt.Sprintf("rating_avg %s", strings.ToUpper(order))
	case "enrolled_count":
		orderClause = fmt.Sprintf("enrolled_count %s", strings.ToUpper(order))
	case "created_at":
		orderClause = fmt.Sprintf("created_at %s", strings.ToUpper(order))
	}

	// Apply pagination and get results
	if err := dbQuery.Order(orderClause).Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, int(total), nil
}

func (cr *DBCourseRepository) GetSearchFilters(query string) (*dto.SearchFilters, error) {
	searchTerm := fmt.Sprintf("%%%s%%", query)

	// Get categories with course count
	var categoryResults []struct {
		Id    uint   `json:"id"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	err := cr.db.Table("course").
		Select("categories.id, categories.name, COUNT(courses.id) as count").
		Joins("JOIN categories ON courses.category_id = categories.id").
		Where("courses.deleted_at IS NULL AND courses.status = ? AND (courses.title ILIKE ? OR courses.description ILIKE ?)",
			"published", searchTerm, searchTerm).
		Group("categories.id, categories.name").
		Scan(&categoryResults).Error

	if err != nil {
		return nil, err
	}

	// Convert to filter options
	categories := make([]dto.FilterOption, len(categoryResults))
	for i, cat := range categoryResults {
		categories[i] = dto.FilterOption{
			Value: fmt.Sprintf("%d", cat.Id),
			Label: cat.Name,
			Count: cat.Count,
		}
	}

	// Static filters (in production, these could be dynamic too)
	levels := []dto.FilterOption{
		{Value: "beginner", Label: "Beginner", Count: 0},
		{Value: "intermediate", Label: "Intermediate", Count: 0},
		{Value: "advanced", Label: "Advanced", Count: 0},
	}

	priceRanges := []dto.FilterOption{
		{Value: "0-50", Label: "Under $50", Count: 0},
		{Value: "50-100", Label: "$50 - $100", Count: 0},
		{Value: "100-200", Label: "$100 - $200", Count: 0},
		{Value: "200+", Label: "Over $200", Count: 0},
	}

	languages := []dto.FilterOption{
		{Value: "vi", Label: "Vietnamese", Count: 0},
		{Value: "en", Label: "English", Count: 0},
	}

	return &dto.SearchFilters{
		Categories:  categories,
		Levels:      levels,
		PriceRanges: priceRanges,
		Languages:   languages,
	}, nil
}

func (cr *DBCourseRepository) GetFeaturedCourses(limit int, filters map[string]interface{}) ([]models.Course, int, error) {
	var courses []models.Course
	var total int64

	query := cr.db.Model(&models.Course{}).
		Preload("Instructor").
		Preload("Category").
		Where("deleted_at IS NULL AND status = ? AND is_featured = ?", "published", true)

	// Apply filters
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get featured courses ordered by rating and enrolled count
	if err := query.
		Order("rating_avg DESC, enrolled_count DESC, created_at DESC").
		Limit(limit).
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, int(total), nil
}

func (cr *DBCourseRepository) FindBySlug(slug string) (*models.Course, error) {
	var course models.Course
	if err := cr.db.
		Preload("Instructor").
		Preload("Category").
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (cr *DBCourseRepository) FindById(courseId uint) (*models.Course, error) {
	var course models.Course
	if err := cr.db.Preload("Instructor").Preload("Category").
		Where("id = ?", courseId).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (cr *DBCourseRepository) UpdateCourseStatus(courseId uint, status string) error {
	return cr.db.Model(&models.Course{}).
		Where("id = ?", courseId).
		Update("status", status).Error
}
