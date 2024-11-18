package user_service

import "time"

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
	Data    interface{} `json:"data,omitempty"` // Use `omitempty` to exclude data if it's nil
}
