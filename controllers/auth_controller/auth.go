package auth_controller

import (
	"errors"
	"fashora-backend/models"
	"fashora-backend/services/auth_service"
	"fashora-backend/services/user_service"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func CheckPhoneNumberExists(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, user_service.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    nil,
		})
		return
	}

	var user models.Users
	err := models.DB.Where("phone = ?", input.PhoneNumber).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, user_service.Response{
				Success: true,
				Status:  http.StatusOK,
				Message: "Phone number does not exist",
				Data:    gin.H{"exists": false},
			})
		} else {
			c.JSON(http.StatusInternalServerError, user_service.Response{
				Success: false,
				Status:  http.StatusInternalServerError,
				Message: "Database error",
				Data:    nil,
			})
		}
	} else {
		c.JSON(http.StatusOK, user_service.Response{
			Success: true,
			Status:  http.StatusOK,
			Message: "Phone number exists",
			Data:    gin.H{"exists": true},
		})
	}
}

func Register(c *gin.Context) {
	var input struct {
		PhoneNumber string     `json:"phone_number" binding:"required"`
		Password    string     `json:"password" binding:"required"`
		UserName    *string    `json:"user_name"`
		Birthday    *time.Time `json:"birthday"`
		Address     *string    `json:"address"`
		DeviceID    *string    `json:"device_id"`
		Gender      *int       `json:"gender"` // should be 0, 1, or 2
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, user_service.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    nil,
		})
		return
	}
	if input.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, user_service.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Phone number is required",
			Data:    nil,
		})
		return
	}

	if !utils.ValidatePhoneNumber(input.PhoneNumber) {
		c.JSON(http.StatusBadRequest, user_service.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Phone number is wrong",
			Data:    nil,
		})
		return
	}

	userWithToken, err := auth_service.Register(user_service.UserInfo(input))
	if err != nil {
		if err.Error() == "phone number already registered" {
			c.JSON(http.StatusBadRequest, user_service.Response{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Phone number already registered",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusBadRequest, user_service.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	responseData := map[string]interface{}{
		"token": userWithToken.Token,
		"user":  userWithToken.User,
	}

	c.JSON(http.StatusCreated, user_service.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    responseData,
	})
}

func Login(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, user_service.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    nil,
		})
		return
	}

	userWithToken, err := auth_service.Login(input.PhoneNumber, input.Password)

	if err != nil {
		switch err.Error() {
		case "phone number not registered":
			c.JSON(http.StatusUnauthorized, user_service.Response{
				Success: false,
				Status:  http.StatusUnauthorized,
				Message: "Phone number not registered",
				Data:    nil,
			})
			return
		case "invalid password":
			c.JSON(http.StatusUnauthorized, user_service.Response{
				Success: false,
				Status:  http.StatusUnauthorized,
				Message: "Invalid password",
				Data:    nil,
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, user_service.Response{
				Success: false,
				Status:  http.StatusInternalServerError,
				Message: "An unexpected error occurred",
				Data:    nil,
			})
			return
		}
	}

	responseData := map[string]interface{}{
		"token": userWithToken.Token,
		"user":  userWithToken.User,
	}

	c.JSON(http.StatusOK, user_service.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Login successful",
		Data:    responseData,
	})
}

func UpdateUser(c *gin.Context) {
	var input user_service.UserInfo

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, user_service.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    nil,
		})
		return
	}

	userInterface, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, user_service.Response{
			Success: false,
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
			Data:    nil,
		})
		return
	}

	user, ok := userInterface.(models.Users)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"status":  http.StatusUnauthorized,
			"message": "Invalid user type",
		})
		return
	}

	if user.Phone != input.PhoneNumber {
		c.JSON(http.StatusUnauthorized, user_service.Response{
			Success: false,
			Status:  http.StatusUnauthorized,
			Message: "Invalid Token",
			Data:    nil,
		})
		return
	}

	if err := user_service.UpdateUserByPhoneNumber(input); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, user_service.Response{
				Success: false,
				Status:  http.StatusNotFound,
				Message: "User does not exist",
				Data:    nil,
			})
			return
		}

	}

	c.JSON(http.StatusOK, user_service.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data:    nil,
	})
}
