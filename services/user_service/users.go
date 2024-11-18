package user_service

import (
	"errors"
	"fashora-backend/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByPhoneNumber(phoneNumber string) (*models.Users, error) {
	var user models.Users
	result := models.DB.Where("phone = ?", phoneNumber).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func CreateNewUser(userInfo UserInfo) (*models.Users, error) {
	var existingUser models.Users
	if err := models.DB.Where("phone = ?", userInfo.PhoneNumber).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create the user instance

	user := models.Users{
		Phone:        userInfo.PhoneNumber,
		PasswordHash: string(hashedPassword),
		UserName:     userInfo.UserName,
		Birthday:     userInfo.Birthday,
		Address:      userInfo.Address,
		DeviceID:     userInfo.DeviceID,
		Gender:       userInfo.Gender,
	}
	if err := models.DB.Create(&user).Error; err != nil {
		return nil, nil
	}

	return &user, nil
}

func UpdateUserByPhoneNumber(userInfoUpdate UserInfo) error {
	updateFields := map[string]interface{}{}
	fmt.Println(userInfoUpdate)

	// Hash password if provided
	if userInfoUpdate.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfoUpdate.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		updateFields["password_hash"] = string(hashedPassword)
	}

	if userInfoUpdate.UserName != nil {
		updateFields["user_name"] = *userInfoUpdate.UserName
	}
	if userInfoUpdate.Birthday != nil {
		updateFields["birthday"] = *userInfoUpdate.Birthday
	}
	if userInfoUpdate.Address != nil {
		updateFields["address"] = *userInfoUpdate.Address
	}
	if userInfoUpdate.DeviceID != nil {
		updateFields["device_id"] = *userInfoUpdate.DeviceID
	}
	if userInfoUpdate.Gender != nil {
		updateFields["gender"] = *userInfoUpdate.Gender
	}

	// If no fields are provided to update, return an error
	if len(updateFields) == 0 {
		return errors.New("no fields to update")
	}

	// Update user record based on phone number
	result := models.DB.Model(models.Users{}).Where("phone = ?", userInfoUpdate.PhoneNumber).Updates(updateFields)
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected (if no rows were affected, it might indicate the user wasn't found)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
