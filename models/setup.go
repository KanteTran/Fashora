package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"login-system/utils"
)

var DB *gorm.DB

func ConnectDatabase() {
	//LoadConfig
	// Use the configuration from AppConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		utils.AppConfig.PostgresHost,
		utils.AppConfig.PostgresUser,
		utils.AppConfig.PostgresPassword,
		utils.AppConfig.PostgresDB,
		utils.AppConfig.PostgresPort,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	database.AutoMigrate(&User_phone{})
	DB = database
}
