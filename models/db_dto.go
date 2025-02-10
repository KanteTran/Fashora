package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"gorm.io/gorm"
)

type Users struct {
	Id           string     `gorm:"primaryKey;size:255"`
	Phone        string     `gorm:"size:255;unique"`
	PasswordHash string     `gorm:"size:255"`
	UserName     *string    `gorm:"size:100"`
	Birthday     *time.Time `gorm:"type:date"`
	Address      *string    `gorm:"size:255"`
	DeviceID     *string    `gorm:"size:100"`
	Gender       *int       `gorm:"check:gender IN (0, 1, 2)"` // Gender: 0 (male), 1 (female), 2 (other), optional
	// Thêm các trường mới
	Height   *float64 `gorm:"type:decimal(5,2)"` // Chiều cao (đơn vị: cm)
	Weight   *float64 `gorm:"type:decimal(5,2)"` // Cân nặng (đơn vị: kg)
	SkinTone *string  `gorm:"size:50;check:skin_tone IN ('light', 'medium', 'dark')"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
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

	// New fields
	Height   *float64 // Chiều cao (đơn vị: cm)
	Weight   *float64 // Cân nặng (đơn vị: kg)
	SkinTone *string  // Màu da (ví dụ: light, medium, dark)
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
	ID          int           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID     string        `json:"store_id" gorm:"not null"`
	Name        string        `json:"name" gorm:"not null"`
	URL         string        `json:"url"`
	ImageURL    string        `json:"image" `
	Description string        `json:"description"`
	Tags        pq.Int64Array `json:"tags" gorm:"type:integer[]"` // Dùng pq.Int64Array để lưu mảng int
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
