package repository

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DBOrderRepository struct {
	db *gorm.DB
}

func NewDBOrderRepository(db *gorm.DB) OrderRepository {
	return &DBOrderRepository{
		db: db,
	}
}

func (or *DBOrderRepository) Create(order *models.Order) error {
	return or.db.Create(order).Error
}

func (or *DBOrderRepository) Update(order *models.Order) error {
	return or.db.Save(order).Error
}

func (or *DBOrderRepository) FindById(orderId uint) (*models.Order, error) {
	var order models.Order
	if err := or.db.Where("id = ? AND deleted_at IS NULL", orderId).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *DBOrderRepository) FindByOrderCode(orderCode string) (*models.Order, error) {
	var order models.Order
	if err := or.db.Where("order_code = ? AND deleted_at IS NULL", orderCode).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *DBOrderRepository) GetUsersOrders(userId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Order, int, error) {
	var orders []models.Order
	var total int64

	query := or.db.Model(&models.Order{}).
		Where("user_id = ? AND deleted_at IS NULL", userId)

	// Apply filters
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Count total
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
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func (or *DBOrderRepository) FindPendingOrderByUserAndCourse(userId, courseId uint) (*models.Order, error) {
	var order models.Order
	err := or.db.Where("user_id = ? AND course_id = ? AND payment_status = ? AND deleted_at IS NULL",
		userId, courseId, "pending").
		First(&order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (or *DBOrderRepository) UpdatePaymentStatus(orderId uint, status string) error {
	updates := map[string]interface{}{
		"payment_status": status,
	}

	if status == "paid" {
		updates["paid_at"] = gorm.Expr("NOW()")
	}

	return or.db.Model(&models.Order{}).
		Where("id = ?", orderId).
		Updates(updates).Error
}

func (or *DBOrderRepository) GetAllOrders(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Order, int, error) {
	var orders []models.Order
	var total int64

	query := or.db.Model(&models.Order{}).
		Preload("User").
		Preload("Course").
		Preload("Course.Instructor").
		Where("orders.deleted_at IS NULL")

	// Apply filters
	for field, value := range filters {
		switch field {
		case "user_id":
			query = query.Where("orders.user_id = ?", value)
		case "course_id":
			query = query.Where("orders.course_id = ?", value)
		case "payment_status":
			query = query.Where("orders.payment_status = ?", value)
		case "payment_method":
			query = query.Where("orders.payment_method = ?", value)
		case "search":
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Joins("LEFT JOIN users ON users.id = orders.user_id").
				Joins("LEFT JOIN courses ON courses.id = orders.course_id").
				Where("orders.order_code ILIKE ? OR users.username ILIKE ? OR users.email ILIKE ? OR courses.title ILIKE ?",
					searchTerm, searchTerm, searchTerm, searchTerm)
		case "date_from":
			query = query.Where("orders.created_at >= ?", value)
		case "date_to":
			query = query.Where("orders.created_at <= ?", value)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("orders.%s %s", orderBy, strings.ToUpper(sortBy)))
	} else {
		query = query.Order("orders.created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func (or *DBOrderRepository) GetOrderStatistics(filters map[string]interface{}) (*dto.OrderStatistics, error) {
	stats := &dto.OrderStatistics{}

	query := or.db.Model(&models.Order{}).Where("deleted_at IS NULL")

	// Apply filters
	for field, value := range filters {
		switch field {
		case "user_id":
			query = query.Where("user_id = ?", value)
		case "course_id":
			query = query.Where("course_id = ?", value)
		case "date_from":
			query = query.Where("created_at >= ?", value)
		case "date_to":
			query = query.Where("created_at <= ?", value)
		}
	}

	// Total orders
	var totalOrders int64
	if err := query.Count(&totalOrders).Error; err != nil {
		return nil, err
	}
	stats.TotalOrders = int(totalOrders)

	// Total revenue (only paid orders)
	var totalRevenue float64
	if err := query.Where("payment_status = ?", "paid").
		Select("COALESCE(SUM(final_price), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}
	stats.TotalRevenue = totalRevenue

	// Count by status
	var pendingCount, completedCount, failedCount, cancelledCount int64

	query.Where("payment_status = ?", "pending").Count(&pendingCount)
	stats.PendingOrders = int(pendingCount)

	query.Where("payment_status = ?", "paid").Count(&completedCount)
	stats.CompletedOrders = int(completedCount)

	query.Where("payment_status = ?", "failed").Count(&failedCount)
	stats.FailedOrders = int(failedCount)

	query.Where("payment_status = ?", "cancelled").Count(&cancelledCount)
	stats.CancelledOrders = int(cancelledCount)

	// Average order value
	if stats.CompletedOrders > 0 {
		stats.AverageOrderValue = totalRevenue / float64(stats.CompletedOrders)
	}

	return stats, nil
}

func (or *DBOrderRepository) UpdateOrderStatus(orderId uint, status string) error {
	updates := map[string]interface{}{
		"payment_status": status,
	}

	if status == "paid" {
		updates["paid_at"] = time.Now()
	}

	return or.db.Model(&models.Order{}).
		Where("id = ?", orderId).
		Updates(updates).Error
}
