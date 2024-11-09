package models

import (
	"time"
)

type Token struct {
	PhoneID     string    `gorm:"primaryKey"`     // Foreign key to user
	Token       string    `gorm:"uniqueIndex"`    // JWT token
	CreatedTime time.Time `gorm:"autoCreateTime"` // When the token was created
	ExpiredTime time.Time // When the token will expire
}
