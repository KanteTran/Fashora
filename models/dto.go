package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (u *Users) BeforeCreate(*gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return
}

type Stores struct {
	Id          string `json:"id" gorm:"primaryKey"`
	Phone       string `json:"phone" gorm:"unique;not null"`
	StoreName   string `json:"store_name" gorm:"not null"`
	Address     string `json:"address" gorm:"not null"`
	Description string `json:"description"`
	Password    string `json:"password" gorm:"not null"`
	Status      int    `json:"status" gorm:"not null"`
	UrlImage    string `json:"url_image" gorm:"not null"`
}

func (u *Stores) BeforeCreate(*gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return
}

type Item struct {
	ID          int      `json:"id" gorm:"primaryKey"`
	StoreID     int      `json:"store_id" gorm:"not null"`
	Name        string   `json:"name" gorm:"not null"`
	URL         string   `json:"url"`
	ImageURLs   []string `json:"image_urls" gorm:"type:text[]"`
	ProductCode string   `json:"product_code"`
}

type ImageRequest struct {
	Image1 []byte `json:"image1"` // First image in binary form
	Image2 []byte `json:"image2"` // Second image in binary form
	Image3 []byte `json:"image3"` // Third image in binary form
}

type ImageResponse struct {
	Status   string `json:"status"`
	ImageURL string `json:"image_url"` // Processed image URL or data
}

type ServiceAccount struct {
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}
