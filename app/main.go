package main

import (
	"fashora-backend/config"
	"fashora-backend/controllers/auth_controller"
	"fashora-backend/controllers/image_controller"
	"fashora-backend/controllers/store_controller"
	"fashora-backend/controllers/try_on_controller"
	"fashora-backend/middlewares"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func setupPublicRoutes(r *gin.Engine) {
	// Auth APIs
	r.POST("/auth/register", auth_controller.Register)
	r.POST("/auth/login", auth_controller.Login)
	r.POST("/auth/check_phone", auth_controller.CheckPhoneNumberExists)

	// Image APIs
	r.POST("/image/push", image_controller.UploadImage)
	r.GET("/image/get", image_controller.GetImageURL)

	// Store APIs
	r.GET("/stores", external.HomePage)
	r.GET("/stores/create-store", external.CreateStorePage)
	r.POST("/stores/create-store", store_controller.CreateStore)
	r.GET("/stores/list-all-store", store_controller.ListStores)
	r.GET("/stores/get_all_items_store", store_controller.GetStoreItemsById)
	r.GET("/stores/get_only_items", store_controller.GetItemsById)

	r.GET("/stores/add-item", store_controller.AddItemPage)
	r.POST("/stores/add-item", store_controller.AddItem)

	// Try On APIs
	r.POST("/try_on/push", try_on_controller.UploadImages)
}

func setupProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Auth APIs requiring authentication
		protected.POST("/auth/update", auth_controller.UpdateUser)

		// Add more authenticated routes here if needed
	}
}

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

	setupPublicRoutes(r)
	setupProtectedRoutes(r)

	// Start server
	go func() {
		err := r.Run(fmt.Sprintf("%s:%s", config.AppConfig.HostServer, config.AppConfig.PortServer))
		if err != nil {
			log.Fatalf("Failed to start REST API server: %v", err)
		}
	}()

	select {}
}
