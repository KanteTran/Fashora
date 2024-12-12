package main

import (
	"fashora-backend/config"
	"fashora-backend/middlewares"
	"fashora-backend/models"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	config.LoadConfig()
	models.ConnectDatabase()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},   // Allow methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow headers
		AllowCredentials: true,                                                // Allow credentials
		MaxAge:           12 * time.Hour,                                      // Max age for preflight requests
	}))

	r.Use(middlewares.SetupApiDocsMiddleware())

	middlewares.SetupPublicRoutes(r)
	middlewares.SetupProtectedRoutes(r)

	// Start server
	go func() {
		err := r.Run(fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port))
		if err != nil {
			log.Fatalf("Failed to start REST API server: %v", err)
		}
	}()

	select {}
}
