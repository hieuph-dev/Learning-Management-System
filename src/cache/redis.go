package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"lms/src/utils"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis khởi tạo kết nối Redis
func InitRedis() error {
	host := utils.GetEnv("REDIS_HOST", "redis")
	port := utils.GetEnv("REDIS_PORT", "6379")
	password := utils.GetEnv("REDIS_PASSWORD", "")
	db := 0 // Database mặc định

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10, // Số lượng kết nối trong pool
		MinIdleConns: 5,  // Số kết nối rảnh tối thiểu
	})

	// Test kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("không thể kết nối Redis: %w", err)
	}

	log.Println("✅ Kết nối Redis thành công!")
	return nil
}

// Set lưu data vào Redis với TTL (Time To Live)
func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Convert value sang JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("không thể marshal data: %w", err)
	}

	// Lưu vào Redis
	return RedisClient.Set(ctx, key, jsonData, ttl).Err()
}

// Get lấy data từ Redis
func Get(ctx context.Context, key string, dest interface{}) error {
	// Lấy data từ Redis
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	// Convert JSON về struct
	return json.Unmarshal([]byte(val), dest)
}

// Delete xóa một key
func Delete(ctx context.Context, keys ...string) error {
	return RedisClient.Del(ctx, keys...).Err()
}

// DeletePattern xóa nhiều keys theo pattern
// Ví dụ: DeletePattern("courses:*") sẽ xóa tất cả keys bắt đầu bằng "courses:"
func DeletePattern(ctx context.Context, pattern string) error {
	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists kiểm tra key có tồn tại không
func Exists(ctx context.Context, key string) bool {
	result, _ := RedisClient.Exists(ctx, key).Result()
	return result > 0
}

// TTL lấy thời gian còn lại của key
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return RedisClient.TTL(ctx, key).Result()
}

// CloseRedis đóng kết nối Redis
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
