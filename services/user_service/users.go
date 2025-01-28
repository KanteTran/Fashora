package user_service

import (
	"context"
	"errors"
	"fashora-backend/database"
	"fashora-backend/logger"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"

	"fashora-backend/models"
)

func GetUserByPhoneNumber(phoneNumber string) (*models.Users, error) {
	var user models.Users
	result := database.GetDBInstance().DB().Where("phone = ?", phoneNumber).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func CreateNewUser(userInfo models.UserInfo) (*models.Users, error) {
	var existingUser models.Users
	if err := database.GetDBInstance().DB().Where(
		"phone = ?", userInfo.PhoneNumber).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := models.Users{
		Phone:        userInfo.PhoneNumber,
		PasswordHash: string(hashedPassword),
		UserName:     userInfo.UserName,
		Birthday:     userInfo.Birthday,
		Address:      userInfo.Address,
		DeviceID:     userInfo.DeviceID,
		Gender:       userInfo.Gender,
	}
	logger.Infof(userInfo.PhoneNumber)
	db := database.GetDBInstance().DB()
	if db == nil {
		logger.Error("Database instance is nil")
		return nil, errors.New("database instance is nil")
	}

	if err := database.GetDBInstance().DB().Create(&user).Error; err != nil {
		logger.Error(fmt.Sprintf("Failed to create user: %v", err))
		return nil, nil
	}

	return &user, nil
}

func UpdateUserByPhoneNumber(userInfoUpdate models.UserInfo) error {
	updateFields := map[string]interface{}{}

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

	if len(updateFields) == 0 {
		return errors.New("no fields to update")
	}

	result := database.GetDBInstance().DB().Model(models.Users{}).Where(
		"phone = ?", userInfoUpdate.PhoneNumber).Updates(updateFields)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// getVerifiedPhoneNumbers retrieves the list of verified phone numbers from Firebase
func GetVerifiedPhoneNumbers(ctx context.Context) ([]string, error) {
	// Initialize Firebase Auth client
	authClient, err := models.FirebaseApp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase Auth client: %v", err)
	}

	// List all users
	iter := authClient.Users(ctx, "")
	var verifiedPhones []string

	for {
		user, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break // Exit the loop when there are no more users
		}
		if err != nil {
			return nil, fmt.Errorf("error listing users: %v", err)
		}

		// Check if the user has a verified phone number
		if user.PhoneNumber != "" {
			verifiedPhones = append(verifiedPhones, user.PhoneNumber)
		}
	}

	return verifiedPhones, nil
}
