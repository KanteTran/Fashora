package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"fashora-backend/config"
	"fashora-backend/models"
	"fashora-backend/services/user_service"
)

var jwtKey = []byte(config.AppConfig.JWT.Secret)

func GenerateJWT(phoneID string) (accessToken string, refreshToken string, err error) {
	accessTokenExpiredTime, _ := strconv.Atoi(config.AppConfig.JWT.AccessTokenExpiration)
	refreshTokenExpiredTime, _ := strconv.Atoi(config.AppConfig.JWT.RefreshTokenExpiration)

	accessClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessTokenExpiredTime) * time.Minute)),
		Issuer:    "fashora-backend",
		Subject:   phoneID,
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(refreshTokenExpiredTime) * time.Hour)),
		Issuer:    "fashora-backend",
		Subject:   phoneID,
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// VerifyJWT kiểm tra JWT Token và trả về thông tin user
func VerifyJWT(tokenString string) (*models.Users, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	// Kiểm tra lỗi parse token
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			// Trả về lỗi rõ ràng nếu token đã hết hạn
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token expired")
			}
			return nil, fmt.Errorf("invalid token: %v", err)
		}
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	// Kiểm tra claims hợp lệ
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Tìm user bằng phoneID (claims.Subject)
	user, err := user_service.GetUserByPhoneNumber(claims.Subject)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return user, nil
}
