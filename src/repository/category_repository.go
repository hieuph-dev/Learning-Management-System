package repository

import (
	"fmt"
	"lms/src/models"

	"gorm.io/gorm"
)

type DBCategoryRepository struct {
	db *gorm.DB
}

func NewDBCategoryRepository(db *gorm.DB) CategoryRepository {
	return &DBCategoryRepository{
		db: db,
	}
}

func (cr *DBCategoryRepository) GetCategories(filters map[string]interface{}) ([]models.Category, int, error) {
	var categories []models.Category
	var total int64

	query := cr.db.Model(&models.Category{}).Where("deleted_at IS NULL")

	// Apply filters
	for field, value := range filters {
		if field == "search" {
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get categories with order
	if err := query.Order("sort_order ASC, name ASC").Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, int(total), nil
}

func (cr *DBCategoryRepository) FindById(id uint) (*models.Category, error) {
	var category models.Category
	if err := cr.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (cr *DBCategoryRepository) Create(category *models.Category) error {
	return cr.db.Create(category).Error
}

func (cr *DBCategoryRepository) FindBySlug(slug string) (*models.Category, bool) {
	var category models.Category
	if err := cr.db.Where("slug = ?", slug).First(&category).Error; err != nil {
		return nil, false
	}
	return &category, true
}

func (cr *DBCategoryRepository) Update(id uint, updates map[string]interface{}) error {
	return cr.db.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error
}

func (cr *DBCategoryRepository) Delete(id uint) error {
	return cr.db.Delete(&models.Category{}, id).Error
}

func (cr *DBCategoryRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := cr.db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (cr *DBCategoryRepository) FindBySlugExcept(slug string, excludeId uint) (*models.Category, bool) {
	var category models.Category
	if err := cr.db.Where("slug = ? AND id != ?", slug, excludeId).First(&category).Error; err != nil {
		return nil, false
	}
	return &category, true
}
