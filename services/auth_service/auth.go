package auth_service

import (
	"errors"
	"fashora-backend/models"
	"fashora-backend/services/user_service"
	"fashora-backend/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserWithToken struct {
	Token string
	User  models.User
}

func Register(userInfo user_service.UserRegisterInfo) (*UserWithToken, error) {
	user, err := user_service.CreateNewUser(userInfo)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return nil, errors.New("phone number already registered")
		}

		return nil, err
	}

	tokenString, err := utils.GenerateJWT(userInfo.PhoneNumber)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &UserWithToken{
		User:  *user,
		Token: tokenString,
	}, nil
}

func Login(phoneNumber string, password string) (*UserWithToken, error) {
	// Find the user by phone number
	user, err := user_service.GetUserByPhoneNumber(phoneNumber)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("phone number not registered")
		}

		return nil, err
	}

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	tokenString, err := utils.GenerateJWT(phoneNumber)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &UserWithToken{
		User:  *user,
		Token: tokenString,
	}, nil
}
