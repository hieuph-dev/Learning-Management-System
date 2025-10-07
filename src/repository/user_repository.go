package repository

import (
	"fmt"
	"lms/src/models"
	"strings"

	"gorm.io/gorm"
)

type DBUserRepository struct {
	db *gorm.DB
}

func NewDBUserRepository(db *gorm.DB) UserRepository {
	return &DBUserRepository{
		db: db,
	}
}

func (ur *DBUserRepository) Create(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *DBUserRepository) FindByEmail(email string) (*models.User, bool) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, false
	}

	return &user, true
}

func (ur *DBUserRepository) FindByUsername(username string) (*models.User, bool) {
	var user models.User
	if err := ur.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, false
	}

	return &user, true
}

func (ur *DBUserRepository) FindById(id uint) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *DBUserRepository) UpdatePassword(userId uint, hashedPassword string) error {
	return ur.db.Model(&models.User{}).Where("id = ?", userId).Update("password", hashedPassword).Error
}

func (ur *DBUserRepository) UpdateProfile(userId uint, updates map[string]interface{}) error {
	return ur.db.Model(&models.User{}).Where("id = ?", userId).Updates(updates).Error
}

func (ur *DBUserRepository) ChangePassword(userId uint, hashedPassword string) error {
	return ur.db.Model(&models.User{}).Where("id = ?", userId).Update("password", hashedPassword).Error
}

func (ur *DBUserRepository) UpdateAvatar(userId uint, avatarURL string) error {
	return ur.db.Model(&models.User{}).Where("id = ?", userId).Update("avatar_url", avatarURL).Error
}

func (ur *DBUserRepository) GetUsersWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.User, int, error) {
	var users []models.User
	var total int64

	query := ur.db.Model(&models.User{})

	// Apply filters
	for field, value := range filters {
		if field == "search" {
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Where("username ILIKE ? OR email ILIKE ? OR full_name ILIKE ?", searchTerm, searchTerm, searchTerm)
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
		query = query.Order((fmt.Sprintf("%s %s", orderBy, strings.ToUpper(sortBy))))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (ur *DBUserRepository) DeleteUser(userId uint) error {
	return ur.db.Delete(&models.User{}, userId).Error
}
