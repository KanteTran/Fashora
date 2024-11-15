package utils

import (
	"fashora-backend/config"
	"fashora-backend/models"
	"fashora-backend/services/user_service"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var jwtKey = []byte(config.AppConfig.JWTSecret) // Replace with a secure key

func GenerateJWT(phoneID string) (string, error) {
	// Generate a new token with configured expiration time
	expired_time, _ := strconv.Atoi(config.AppConfig.JwtExpirationHours)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expired_time) * time.Hour)),
		Issuer:    "fashora-backend",
		Subject:   phoneID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Kiểm tra token có hết hạn không
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, jwt.NewValidationError("Token expired", jwt.ValidationErrorExpired)
		}

		// Kiểm tra token có tồn tại trong database không
		user, err := user_service.GetUserByPhoneNumber(claims.Subject)
		if user != nil {
			return user, nil
		}

		// user not exist
		if err == gorm.ErrRecordNotFound {
			return nil, jwt.NewValidationError("Invalid token", jwt.ValidationErrorMalformed)
		}

		// other db errors
		return nil, err
	}

	// claims not ok or token not valid
	return nil, jwt.NewValidationError("Invalid token", jwt.ValidationErrorMalformed)
}
