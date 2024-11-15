package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"fashora-backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("AuthMiddleware is called")

		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// Trả về lỗi nếu không có token
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Xóa tiền tố "Bearer " nếu có
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Kiểm tra tính hợp lệ của token bằng VerifyJWT
		user, err := utils.VerifyJWT(tokenString)
		if err != nil {
			// Trả về lỗi nếu token không hợp lệ hoặc hết hạn
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token hợp lệ, lưu thông tin user ID vào context
		c.Set("user", user)
		c.Next()
	}
}
