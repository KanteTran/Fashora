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
	// Auth APIs
	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)
	r.POST("/auth/check_phone", auth.CheckPhoneNumberExists)

	// Store APIs
	r.GET("/stores", external.HomePage)
	r.GET("/stores/create-store", external.CreateStorePage)
	r.POST("/stores/create-store", store.CreateStore)
	r.GET("/stores/add-item", store.AddItemPage)
	r.POST("/stores/add-item", store.AddItem)
	r.POST("/stores/get-items-by-tags", recommend.GetItemsByTags)
	// Score APIs
	r.POST("/image/scoring", scoring.ScoreImage)
	//r.POST("image/tag", tagging.TagImage)
	// Recommend API
	r.POST("/recommend/gen_tags", recommend.GenTagRecommend)
	// Get version
	r.GET("/version", external.Version)

}

func SetupProtectedRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		// Auth APIs requiring authentication
		protected.POST("/auth/update", auth.Update)
		protected.POST("/try_on/push", try_on.UploadImages)
		protected.POST("/inventory/add-item", inventory.AddInventory)
		protected.GET("/inventory/all-items", inventory.ListInventories)
		protected.DELETE("/inventory/del-item", inventory.DeleteInventory)
		protected.GET("/stores/list-all-store", store.ListStores)
		protected.GET("/stores/get_all_items_store", store.GetStoreItemsById)
		protected.GET("/stores/get_only_items", store.GetItemsById)
		protected.POST("/try_on/segment", try_on.Segment)

		// Add more authenticated routes here if needed
	}
}
