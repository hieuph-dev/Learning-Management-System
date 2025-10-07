package middleware

import (
	"fmt"
	"lms/src/cache"
	"lms/src/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware tạo middleware cache với TTL tùy chỉnh
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Chỉ cache cho GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Tạo cache key từ URL và query parameters
		cacheKey := generateCacheKey(c)

		// Thử lấy từ cache
		var cachedResponse interface{}
		err := cache.Get(c.Request.Context(), cacheKey, &cachedResponse)

		// Nếu tìm thấy trong cache
		if err == nil {
			c.Header("X-Cache-Status", "HIT") // Đánh dấu là cache hit
			utils.ResponseSuccess(c, http.StatusOK, cachedResponse)
			c.Abort()
			return
		}

		// Nếu không tìm thấy hoặc có lỗi (trừ lỗi key not found)
		if err != redis.Nil {
			// Log lỗi nhưng vẫn tiếp tục xử lý request
			fmt.Printf("⚠️ Redis error: %v\n", err)
		}

		// Tạo custom writer để capture response
		blw := &bodyLogWriter{
			body:           []byte{},
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		c.Header("X-Cache-Status", "MISS") // Đánh dấu là cache miss

		// Xử lý request bình thường
		c.Next()

		// Sau khi xử lý xong, lưu vào cache nếu response thành công
		if c.Writer.Status() == http.StatusOK && len(blw.body) > 0 {
			// Parse response body
			var responseData map[string]interface{}
			if err := utils.ParseJSON(blw.body, &responseData); err == nil {
				// Chỉ cache phần data, bỏ qua các metadata khác
				if data, ok := responseData["data"]; ok {
					if err := cache.Set(c.Request.Context(), cacheKey, data, ttl); err != nil {
						fmt.Printf("⚠️ Không thể lưu vào cache: %v\n", err)
					} else {
						fmt.Printf("✅ Đã cache: %s (TTL: %v)\n", cacheKey, ttl)
					}
				}
			}
		}
	}
}

// generateCacheKey tạo unique cache key
func generateCacheKey(c *gin.Context) string {
	// Kết hợp path và query parameters
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	if query != "" {
		return fmt.Sprintf("cache:%s?%s", path, query)
	}
	return fmt.Sprintf("cache:%s", path)
}

// bodyLogWriter để capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

// InvalidateCachePattern xóa cache theo pattern
func InvalidateCachePattern(pattern string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Xử lý request trước
		c.Next()

		// Sau khi xử lý xong, xóa cache nếu request thành công
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			if err := cache.DeletePattern(c.Request.Context(), pattern); err != nil {
				fmt.Printf("⚠️ Không thể xóa cache pattern %s: %v\n", pattern, err)
			} else {
				fmt.Printf("🗑️ Đã xóa cache: %s\n", pattern)
			}
		}
	}
}
