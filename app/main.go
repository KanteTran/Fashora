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

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	config.LoadConfig()
	models.ConnectDatabase()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},   // Allow methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow headers
		AllowCredentials: true,                                                // Allow credentials
		MaxAge:           12 * time.Hour,                                      // Max age for preflight requests
	}))

	r.POST("/auth/register", auth_controller.Register)
	r.POST("/auth/login", auth_controller.Login)
	r.POST("/auth/check_phone", auth_controller.CheckPhoneNumberExists)

	r.POST("/image/push", image_controller.UploadImage)
	r.GET("/image/get", image_controller.GetImageURL)

	r.GET("/stores", external.HomePage)
	r.GET("/stores/create-store", external.CreateStorePage)
	r.POST("/stores/create-store", store_controller.CreateStore)
	r.GET("/stores/list-all-store", store_controller.ListStores)

	r.GET("/stores/add-item", store_controller.AddItemPage)
	r.POST("/stores/add-item", store_controller.AddItem)

	r.POST("/try_on/push", try_on_controller.UploadImagesV2)
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/auth/update", auth_controller.UpdateUser)
	}

	go func() {
		err := r.Run(fmt.Sprintf("%s:%s", config.AppConfig.HostServer, config.AppConfig.PortServer))
		if err != nil {
			log.Fatalf("Failed to start REST API server: %v", err)
		}
	}()

	select {}

}
