package middlewares

import (
	"github.com/gin-gonic/gin"

	"fashora-backend/handler/auth"
	"fashora-backend/handler/inventory"
	"fashora-backend/handler/recommend"
	"fashora-backend/handler/scoring"
	"fashora-backend/handler/store"
	"fashora-backend/handler/try_on"
	"fashora-backend/services/external"
)

func SetupPublicRoutes(r *gin.Engine) {
	storeHandler := store.NewHandlerStore()

	// Auth APIs
	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)
	r.POST("/auth/check_phone", auth.CheckPhoneNumberExists)

	// Store APIs
	r.GET("/stores", external.HomePage)
	r.GET("/stores/create-store", external.CreateStorePage)
	r.POST("/stores/create-store", storeHandler.CreateStore)
	r.GET("/stores/add-item", storeHandler.AddItemPage)
	r.POST("/stores/add-item", storeHandler.AddItem)
	r.POST("/stores/get-items-by-tags", recommend.GetItemsByTags)

	r.GET("/stores/list-all-store", storeHandler.ListStores)
	r.GET("/stores/get_all_items_store", storeHandler.GetStoreItemsById)
	r.GET("/stores/get_only_items", storeHandler.GetItemsById)
	r.POST("/try_on/segment", try_on.Segment)
	r.POST("/image/scoring", scoring.ScoreImage)
	r.POST("/recommend/gen_tags", recommend.GenTagRecommend)
	r.POST("/try_on/push", try_on.UploadImages)
	// Get version
	r.GET("/version", external.Version)

}

func SetupProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		// Auth APIs requiring authentication
		protected.POST("/auth/update", auth.Update)
		//protected.POST("/try_on/push", try_on.UploadImages)
		protected.POST("/inventory/add-item", inventory.AddInventory)
		protected.GET("/inventory/all-items", inventory.ListInventories)
		protected.DELETE("/inventory/del-item", inventory.DeleteInventory)

		// Add more authenticated routes here if needed
	}
}
