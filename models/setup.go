package models

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"fashora-backend/config"
	"fashora-backend/logger"
)

var DB *gorm.DB
var FirebaseApp *firebase.App

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
		logger.Errorf("Failed to connect to the database: %s", err)
	}

	err = database.AutoMigrate(&Users{}, &Stores{}, &Item{}, &Inventory{})
	if err != nil {
		return
	}
	DB = database
	logger.Infof("Successfully connected to the database")

	// Firebase initialization
	opt := option.WithCredentialsFile(config.AppConfig.FireBase.FileKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf("Failed to initialize Firebase: %s", err)
		return
	}

	FirebaseApp = app
	logger.Infof("Successfully connected to Firebase")
}
