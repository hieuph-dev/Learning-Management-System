package repository

import (
	"lms/src/models"
	"time"

	"gorm.io/gorm"
)

type DBPasswordResetRepository struct {
	db *gorm.DB
}

func NewDBPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &DBPasswordResetRepository{
		db: db,
	}
}

func (pr *DBPasswordResetRepository) Create(reset *models.PasswordReset) error {
	return pr.db.Create(reset).Error
}

func (pr *DBPasswordResetRepository) FindByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := pr.db.Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (pr *DBPasswordResetRepository) MarkAsUsed(id uint) error {
	return pr.db.Model(&models.PasswordReset{}).Where("id = ?", id).Update("used", true).Error
}

func (pr *DBPasswordResetRepository) DeleteExpired() error {
	return pr.db.Where("expires_at < ? OR used = true", time.Now()).Delete(&models.PasswordReset{}).Error
}

func (pr *DBPasswordResetRepository) DeleteByEmail(email string) error {
	return pr.db.Where("email = ?", email).Delete(&models.PasswordReset{}).Error
}
