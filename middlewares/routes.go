package middlewares

import (
	"fashora-backend/controllers/auth_controller"
	"fashora-backend/controllers/inventory_controller"
	"fashora-backend/controllers/store_controller"
	"fashora-backend/controllers/try_on_controller"
	"fashora-backend/services/external"
	"github.com/gin-gonic/gin"
)

func SetupPublicRoutes(r *gin.Engine) {
	// Auth APIs
	r.POST("/auth/register", auth_controller.Register)
	r.POST("/auth/login", auth_controller.Login)
	r.POST("/auth/check_phone", auth_controller.CheckPhoneNumberExists)

	// Store APIs
	r.GET("/stores", external.HomePage)
	r.GET("/stores/create-store", external.CreateStorePage)
	r.POST("/stores/create-store", store_controller.CreateStore)
	r.GET("/stores/list-all-store", store_controller.ListStores)
	r.GET("/stores/get_all_items_store", store_controller.GetStoreItemsById)
	r.GET("/stores/get_only_items", store_controller.GetItemsById)

	r.GET("/stores/add-item", store_controller.AddItemPage)
	r.POST("/stores/add-item", store_controller.AddItem)

}

func SetupProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		// Auth APIs requiring authentication
		protected.POST("/auth/update", auth_controller.Update)
		protected.POST("/try_on/push", try_on_controller.UploadImages)
		protected.POST("/inventory/add-item", inventory_controller.AddInventory)
		protected.GET("/inventory/all-items", inventory_controller.ListInventories)
		protected.DELETE("/inventory/del-item", inventory_controller.DeleteInventory)
		// Add more authenticated routes here if needed
	}
}
