package main

import (
	"fashora-backend/config"
	"fashora-backend/controllers/auth_controller"
	"fashora-backend/controllers/user_controller"
	"fashora-backend/middlewares"
	"fashora-backend/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.LoadConfig()

	models.ConnectDatabase()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},   // Allow methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow headers
		AllowCredentials: true,                                                // Allow credentials
		MaxAge:           12 * time.Hour,                                      // Max age for preflight requests
	}))

	// Routes
	r.POST("/auth/register", auth_controller.Register)
	r.POST("/auth/login", auth_controller.Login)

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.PATCH("/user/update_user", user_controller.UpdateUser)
	}

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
