package utils

import (
	"errors"
	"fashora-backend/config"
	"fashora-backend/models"
	"fashora-backend/services/user_service"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var jwtKey = []byte(config.AppConfig.JWTSecret)

func GenerateJWT(phoneID string) (string, error) {
	expiredTime, _ := strconv.Atoi(config.AppConfig.JwtExpirationHours)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredTime) * time.Hour)),
		Issuer:    "fashora-backend",
		Subject:   phoneID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenString string) (*models.Users, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, jwt.NewValidationError("Token expired", jwt.ValidationErrorExpired)
		}

		user, err := user_service.GetUserByPhoneNumber(claims.Subject)
		if user != nil {
			return user, nil
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, jwt.NewValidationError("Invalid token", jwt.ValidationErrorMalformed)
		}

		return nil, err
	}

	return nil, jwt.NewValidationError("Invalid token", jwt.ValidationErrorMalformed)
}
