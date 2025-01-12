package auth_service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"fashora-backend/models"
	"fashora-backend/services/user_service"
	"fashora-backend/utils"
)

type UserWithToken struct {
	Token string
	User  models.Users
}

func Register(userInfo models.UserInfo) (*UserWithToken, error) {
	user, err := user_service.CreateNewUser(userInfo)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("phone number already registered")
		}

		return nil, err
	}
	//return nil

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
	user, err := user_service.GetUserByPhoneNumber(phoneNumber)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("phone number not registered")
		}

		return nil, err
	}

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
