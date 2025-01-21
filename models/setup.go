package models

import (
	"context"
	"fashora-backend/config"
	"fashora-backend/logger"
	//"fashora-backend/services/external"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

//var GeminiApp *external.GeminiApp

func ConnectDatabase() {
	// Firebase initialization
	opt := option.WithCredentialsFile(config.AppConfig.FireBase.FileKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf("Failed to initialize Firebase: %s", err)
		return
	}

	FirebaseApp = app

	//GeminiApp = external.InitGemini(config.AppConfig.Model.GeminiAPI)
	logger.Infof("Successfully connected to Firebase")
}
