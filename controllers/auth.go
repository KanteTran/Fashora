package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"login-system/models"
	"login-system/utils"
	"net/http"
	"strconv"
	"time"
)

func Register(c *gin.Context) {
	var input struct {
		PhoneNumber string     `json:"phone_number" binding:"required"`
		Password    string     `json:"password" binding:"required"`
		UserName    *string    `json:"user_name"` // Optional field
		Birthday    *time.Time `json:"birthday"`  // Optional field, format: "YYYY-MM-DD"
		Address     *string    `json:"address"`   // Optional field
		DeviceID    *string    `json:"device_id"` // Optional field
		Gender      *int       `json:"gender"`    // Optional field, should be 0, 1, or 2
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	hashedPasswordStr := string(hashedPassword)
	// Create the user instance
	user := models.User_phone{
		PhoneID:      input.PhoneNumber,
		PasswordHash: &hashedPasswordStr,
		UserName:     input.UserName,
		Birthday:     input.Birthday,
		Address:      input.Address,
		DeviceID:     input.DeviceID,
		Gender:       input.Gender,
	}

	// Save the user to the database
	result := models.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already registered"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the user by phone number
	var user models.User_phone
	if err := models.DB.Where("phone_id = ?", input.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		return
	}

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		return
	}

	// Check if there's an existing token for this user that hasn't expired
	var userToken models.Token
	if err := models.DB.Where("phone_id = ?", input.PhoneNumber).First(&userToken).Error; err == nil {
		// Token exists, check if it's still valid
		if time.Now().Before(userToken.ExpiredTime) {
			// Token is still valid, return the existing token
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"token":   userToken.Token,
			})
			return
		}
		// Token is expired, so we'll create a new one
	}

	// Generate a new token with a 24-hour expiration time

	experied_time, err := strconv.Atoi(utils.AppConfig.JwtExpirationHours)
	tokenString, err := utils.GenerateJWT(input.PhoneNumber, time.Duration(experied_time)*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save the new token in the Token table
	userToken = models.Token{
		PhoneID:     input.PhoneNumber,
		Token:       tokenString,
		CreatedTime: time.Now(),
		ExpiredTime: time.Now().Add(time.Duration(experied_time) * time.Hour), // Set expiration time to 72 hours
	}

	// Upsert the token (create if not exist, update if exist)
	if err := models.DB.Save(&userToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// Return the new token to the client
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func Update(c *gin.Context) {
	var input struct {
		PhoneNumber string     `json:"phone_number" binding:"required"`
		Password    *string    `json:"password"`  // Optional, hashed if provided
		UserName    *string    `json:"user_name"` // Optional
		Birthday    *time.Time `json:"birthday"`  // Optional, format: "YYYY-MM-DD"
		Address     *string    `json:"address"`   // Optional
		DeviceID    *string    `json:"device_id"` // Optional
		Gender      *int       `json:"gender"`    // Optional, should be 0, 1, or 2
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the user by phone number
	var user models.User_phone
	if err := models.DB.Where("phone_id = ?", input.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields if provided in input
	if input.Password != nil {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		hashedPasswordStr := string(hashedPassword)
		user.PasswordHash = &hashedPasswordStr
	}
	if input.UserName != nil {
		user.UserName = input.UserName
	}
	if input.Birthday != nil {
		user.Birthday = input.Birthday
	}
	if input.Address != nil {
		user.Address = input.Address
	}
	if input.DeviceID != nil {
		user.DeviceID = input.DeviceID
	}
	if input.Gender != nil {
		user.Gender = input.Gender
	}

	// Manually set UpdatedAt to current time (optional, GORM does this automatically on Save)
	user.UpdatedAt = time.Now()

	// Save the changes to the database
	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
