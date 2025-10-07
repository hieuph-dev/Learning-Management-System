package db

import (
	"context"
	"fmt"
	"lms/src/config"
	"lms/src/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	connStr := config.NewDBConfig().DNS() // Lấy thông tin kết nối DB từ config

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Bật logging để xem các câu SQL được thực thi
	}

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{ // Tạo kết nối với cấu hình đã định nghĩa
		DSN: connStr,
	}), config)

	if err != nil {
		return fmt.Errorf("error opening DB connection: %w", err)
	}

	sqlDB, err := DB.DB() // Lấy underlying sql.DB
	if err != nil {
		return fmt.Errorf("error getting sql.DB: %v", err)
	}

	// Connection pool

	sqlDB.SetMaxOpenConns(50) // Tối đa 50 kết nối đồng thời

	sqlDB.SetMaxIdleConns(10) // Giữ 10 kết nối rảnh

	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Kết nối sống tối đa 30 phút

	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // Kết nối rảnh tối đa 5 phút

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Ping database với timeout 5 giây
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close() // Đóng kết nối nếu ping thất bại
		return fmt.Errorf("DB ping error: %w", err)
	}

	err = DB.AutoMigrate( // Tự động tạo/cập nhật bảng dựa trên struct
		&models.User{},
		&models.PasswordReset{},
		&models.Category{},
		&models.Course{},
		&models.Lesson{},
		&models.Enrollment{},
		&models.Progress{},
		&models.Review{},
		&models.Coupon{},
		&models.Order{},
	)

	if err != nil {
		sqlDB.Close()
		return fmt.Errorf("error running migration: %w", err)
	}

	log.Println("Connected and migrated successfully")

	return nil
}
