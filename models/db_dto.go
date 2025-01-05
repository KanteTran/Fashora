package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	Id           string         `gorm:"primaryKey;size:255"`
	Phone        string         `gorm:"size:255;unique"`
	PasswordHash string         `gorm:"size:255"`
	UserName     *string        `gorm:"size:100"`
	Birthday     *time.Time     `gorm:"type:date"`
	Address      *string        `gorm:"size:255"`
	DeviceID     *string        `gorm:"size:100"`
	Gender       *int           `gorm:"check:gender IN (0, 1, 2)"` // Gender: 0 (male), 1 (female), 2 (other), optional
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Inventory struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	StoreID   string         `gorm:"size:255" json:"store_id"`
	ItemID    string         `gorm:"size:255" json:"item_id"`
	Name      string         `gorm:"size:255" json:"name"`
	URL       string         `gorm:"size:255" json:"url"`
	ImageURL  string         `gorm:"size:2555" json:"image_url"`
	UserID    string         `gorm:"size:255" json:"user_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *Users) BeforeCreate(*gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return
}

type UserInfo struct {
	PhoneNumber string
	Password    string
	UserName    *string
	Birthday    *time.Time
	Address     *string
	DeviceID    *string
	Gender      *int
}

type Response struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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
	Type        string `json:"type" gorm:"not null;default:'1'"`
}

func (u *Stores) BeforeCreate(*gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return
}

type Item struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID     string `json:"store_id" gorm:"not null"`
	Name        string `json:"name" gorm:"not null"`
	URL         string `json:"url"`
	ImageURL    string `json:"image" `
	Description string `json:"description"`
}

type ImageRequest struct {
	Image1 []byte `json:"image1"`
	Image2 []byte `json:"image2"`
	Image3 []byte `json:"image3"`
}

type ImageResponse struct {
	Status   string `json:"status"`
	ImageURL string `json:"image_url"`
}

type Image struct {
	FormKey    string
	BucketName string
}

type UserWithToken struct {
	Token string
	User  Users
}
