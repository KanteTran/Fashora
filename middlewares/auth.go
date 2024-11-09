package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"login-system/models"
	"login-system/utils"
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
		_, err := utils.VerifyJWT(tokenString)
		if err != nil {
			// Trả về lỗi nếu token không hợp lệ hoặc hết hạn
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Kiểm tra token có tồn tại trong database không
		var userToken models.Token
		if err := models.DB.Where("token = ?", tokenString).First(&userToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Kiểm tra token có hết hạn không
		if time.Now().After(userToken.ExpiredTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired, please log in again"})
			c.Abort()
			return
		}

		// Token hợp lệ, lưu thông tin user ID vào context
		c.Set("userID", userToken.PhoneID)
		c.Next()
	}
}
