package middleware

import (
	"lms/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy token từ header Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  utils.ErrCodeUnauthorized,
			})
			// Dừng ngay lập tức
			ctx.Abort()
			return
		}

		// Kiểm tra format: "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
				"code":  utils.ErrCodeUnauthorized,
			})
			ctx.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  utils.ErrCodeUnauthorized,
			})
			ctx.Abort()
			return
		}

		// Kiểm tra xem token có phải là access token không
		if claims.Subject != "access" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token type",
				"code":  utils.ErrCodeUnauthorized,
			})
			ctx.Abort()
			return
		}

		// Lưu thông tin User vào Context
		ctx.Set("user_id", claims.UserId)
		ctx.Set("username", claims.Username)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}
