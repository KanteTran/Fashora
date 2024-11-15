package user_service

import (
	"errors"
	"fashora-backend/models"

	"golang.org/x/crypto/bcrypt"
)

func GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	result := models.DB.First(&user, phoneNumber)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func CreateNewUser(userInfo UserRegisterInfo) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create the user instance
	user := models.User{
		PhoneID:      userInfo.PhoneNumber,
		PasswordHash: string(hashedPassword),
		UserName:     userInfo.UserName,
		Birthday:     userInfo.Birthday,
		Address:      userInfo.Address,
		DeviceID:     userInfo.DeviceID,
		Gender:       userInfo.Gender,
	}

	// Save the user to the database
	result := models.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// TODO phone is used as primary key, so what if user wants to update their phone number?
func UpdateUserByPhoneNumber(userInfo UserUpdateInfo) error {
	updateFields := map[string]interface{}{}

	// Conditionally add fields to update
	if userInfo.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*userInfo.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		updateFields["PasswordHash"] = string(hashedPassword)
	}
	if userInfo.UserName != nil {
		updateFields["UserName"] = *userInfo.UserName
	}
	if userInfo.Birthday != nil {
		updateFields["Birthday"] = *userInfo.Birthday
	}
	if userInfo.Address != nil {
		updateFields["Address"] = *userInfo.Address
	}
	if userInfo.DeviceID != nil {
		updateFields["DeviceID"] = *userInfo.DeviceID
	}
	if userInfo.Gender != nil {
		updateFields["Gender"] = *userInfo.Gender
	}

	// update user
	return models.DB.Where("phone_id = ?", userInfo.PhoneNumber).Updates(updateFields).Error
}
