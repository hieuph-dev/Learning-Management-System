package repository

import (
	"lms/src/models"
	"time"

	"gorm.io/gorm"
)

type DBCouponRepository struct {
	db *gorm.DB
}

func NewDBCouponRepository(db *gorm.DB) CouponRepository {
	return &DBCouponRepository{
		db: db,
	}
}

func (cr *DBCouponRepository) FindByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := cr.db.Where("code = ? AND is_active = ? AND deleted_at IS NULL", code, true).
		First(&coupon).Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (cr *DBCouponRepository) FindById(id uint) (*models.Coupon, error) {
	var coupon models.Coupon
	err := cr.db.Where("id = ? AND deleted_at IS NULL", id).
		First(&coupon).Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (cr *DBCouponRepository) IncrementUsedCount(couponId uint) error {
	return cr.db.Model(&models.Coupon{}).
		Where("id = ?", couponId).
		Update("used_count", gorm.Expr("used_count + 1")).Error
}

func (cr *DBCouponRepository) IsValidCoupon(coupon *models.Coupon) bool {
	now := time.Now()

	// Check if coupon is active
	if !coupon.IsActive {
		return false
	}

	// Check valid from date
	if coupon.ValidFrom != nil && now.Before(*coupon.ValidFrom) {
		return false
	}

	// Check valid to date
	if coupon.ValidTo != nil && now.After(*coupon.ValidTo) {
		return false
	}

	// Check usage limit
	if coupon.UsageLimit != nil && coupon.UsedCount >= *coupon.UsageLimit {
		return false
	}

	return true
}

func (cr *DBCouponRepository) GetCouponsWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Coupon, int, error) {
	var coupons []models.Coupon
	var total int64

	query := cr.db.Model(&models.Coupon{}).Where("deleted_at IS NULL")

	// Apply filters
	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if searchCode, ok := filters["search_code"].(string); ok && searchCode != "" {
		query = query.Where("code LIKE ?", "%"+searchCode+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" && sortBy != "" {
		query = query.Order(orderBy + " " + sortBy)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&coupons).Error; err != nil {
		return nil, 0, err
	}

	return coupons, int(total), nil
}

func (cr *DBCouponRepository) Create(coupon *models.Coupon) error {
	return cr.db.Create(coupon).Error
}

func (cr *DBCouponRepository) Update(couponId uint, updates map[string]interface{}) error {
	return cr.db.Model(&models.Coupon{}).
		Where("id = ? AND deleted_at IS NULL", couponId).
		Updates(updates).Error
}

func (cr *DBCouponRepository) Delete(couponId uint) error {
	return cr.db.Where("id = ?", couponId).
		Delete(&models.Coupon{}).Error
}

func (cr *DBCouponRepository) FindByCodeExcept(code string, excludeId uint) (*models.Coupon, bool) {
	var coupon models.Coupon
	err := cr.db.Where("code = ? AND id != ? AND deleted_at IS NULL", code, excludeId).
		First(&coupon).Error

	if err != nil {
		return nil, false
	}

	return &coupon, true
}
