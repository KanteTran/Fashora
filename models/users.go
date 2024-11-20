package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Users struct {
	Id           string         `gorm:"primaryKey;size:255"`       // Primary key, required
	Phone        string         `gorm:"size:255;unique"`           // Primary key, required
	PasswordHash string         `gorm:"size:255"`                  // Hashed password
	UserName     *string        `gorm:"size:100"`                  // Username, optional (nullable)
	Birthday     *time.Time     `gorm:"type:date"`                 // Birthday, optional (nullable)
	Address      *string        `gorm:"size:255"`                  // Address, optional (nullable)
	DeviceID     *string        `gorm:"size:100"`                  // Device identifier, optional (nullable)
	Gender       *int           `gorm:"check:gender IN (0, 1, 2)"` // Gender: 0 (male), 1 (female), 2 (other), optional
	CreatedAt    time.Time      `gorm:"autoCreateTime"`            // Automatically sets time on creation
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`            // Automatically updates time on modification
	DeletedAt    gorm.DeletedAt `gorm:"index"`                     // Soft delete field, optional
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = uuid.New().String() // Generate a new UUID if Id is empty
	}
	return
}
