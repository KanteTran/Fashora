package main

import (
	"github.com/gin-gonic/gin"
	"login-system/controllers"
	"login-system/middlewares"
	"login-system/models"
	"login-system/utils"
)

func main() {
	r := gin.Default()
	utils.LoadConfig()

	// Kết nối database
	models.ConnectDatabase()

	// Routes
	r.POST("/user/register", controllers.Register)
	r.POST("/user/update_user", controllers.Update)
	r.POST("/user/login", controllers.Login)

	protected := r.Group("/user")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.PUT("/user/update_user", controllers.Update) // Route update người dùng
	}

	r.Run(":8080")
}
