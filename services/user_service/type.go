package user_service

import "time"

type UserRegisterInfo struct {
	PhoneNumber string
	Password    string
	UserName    *string
	Birthday    *time.Time
	Address     *string
	DeviceID    *string
	Gender      *int
}

type UserUpdateInfo struct {
	PhoneNumber string
	Password    *string
	UserName    *string
	Birthday    *time.Time
	Address     *string
	DeviceID    *string
	Gender      *int
}
