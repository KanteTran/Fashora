package auth_service

import (
	"errors"
	"fashora-backend/config"
	"fashora-backend/database"
	"fashora-backend/logger"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"fashora-backend/models"
	"fashora-backend/services/user_service"
	"fashora-backend/utils"
)

type UserWithToken struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         models.Users `json:"user"`
}

func Register(userInfo models.UserInfo) (*UserWithToken, error) {
	user, err := user_service.CreateNewUser(userInfo)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("phone number already registered")
		}
		return nil, err
	}

	refreshTokenExpiredTime, _ := strconv.Atoi(config.AppConfig.JWT.RefreshTokenExpiration)
	// Tạo Access Token & Refresh Token
	accessToken, refreshToken, err := utils.GenerateJWT(userInfo.PhoneNumber)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refreshExpiresAt := time.Now().Add(time.Duration(refreshTokenExpiredTime) * time.Hour)
	// Lưu Refresh Token vào Database (hoặc Redis)
	err = SaveRefreshToken(user.Id, refreshToken, refreshExpiresAt)
	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return &UserWithToken{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	refreshTokenExpiredTime, _ := strconv.Atoi(config.AppConfig.JWT.RefreshTokenExpiration)

	// Tạo Access Token & Refresh Token
	accessToken, refreshToken, err := utils.GenerateJWT(phoneNumber)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refreshExpiresAt := time.Now().Add(time.Duration(refreshTokenExpiredTime) * time.Hour)
	// Lưu Refresh Token vào Database (hoặc Redis)
	err = SaveRefreshToken(user.Id, refreshToken, refreshExpiresAt)

	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return &UserWithToken{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func SaveRefreshToken(userID string, refreshToken string, expiresAt time.Time) error {
	db := database.GetDBInstance().DB()
	if db == nil {
		logger.Error("Database instance is nil")
		return errors.New("database instance is nil")
	}

	err := db.Exec(`INSERT INTO refresh_tokens (user_id, token, expires_at)
						VALUES (?, ?, ?)
						ON CONFLICT (user_id) DO UPDATE 
						SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at;`,
		userID, refreshToken, expiresAt,
	).Error

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to save refresh token: %v", err))
		return errors.New("failed to save refresh token")
	}

	return nil
}

func ValidateRefreshToken(refreshToken string) (uint, error) {
	var userID uint
	db := database.GetDBInstance().DB()
	if db == nil {
		logger.Error("Database instance is nil")
		return 0, errors.New("database instance is nil")
	}

	// Kiểm tra token có trong DB không, và còn hạn sử dụng không
	err := db.Raw(`
		SELECT user_id FROM refresh_tokens WHERE token = ? AND expires_at > ?`,
		refreshToken, time.Now(),
	).Scan(&userID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("invalid or expired refresh token")
		}
		logger.Error(fmt.Sprintf("Error validating refresh token: %v", err))
		return 0, errors.New("database error")
	}

	return userID, nil
}
