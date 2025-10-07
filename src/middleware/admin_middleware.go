package middleware

import (
	"lms/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware kiểm tra xem user có role admin không
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy role từ context (đã được set bởi AuthMiddleware)
		userRole, exists := ctx.Get("user_role")
		if !exists {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "User role not found in context",
				"code":  utils.ErrCodeForbidden,
			})
			ctx.Abort()
			return
		}

		// Kiểm tra xem có phải admin không
		if userRole.(string) != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. Admin role required",
				"code":  utils.ErrCodeForbidden,
			})
			ctx.Abort()
			return
		}

		// Cho phép tiếp tục
		ctx.Next()
	}
}
