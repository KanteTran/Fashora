package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fashora-backend/config"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		config.AppConfig.Postgres.Host,
		config.AppConfig.Postgres.User,
		config.AppConfig.Postgres.Password,
		config.AppConfig.Postgres.DB,
		config.AppConfig.Postgres.Port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = database.AutoMigrate(&Users{}, &Stores{}, &Item{}, &Inventory{})
	if err != nil {
		return
	}
	DB = database
	log.Println("Successfully connected to the database")
}
