package interfaces

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *ShopHandler) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	v1 := r.Group("/v1")
	{
		// Shop routes
		shopRoutes := v1.Group("/shop")
		shopRoutes.Use(AuthMiddleware()) // Apply AuthMiddleware to all /shop routes in this group
		{
			shopRoutes.POST("", handler.CreateShop)             // POST /v1/shop
			shopRoutes.GET("/my", handler.GetUserShops)       // GET /v1/shop/my
			shopRoutes.GET("/:shop_id", handler.GetShop)      // GET /v1/shop/:shop_id
			shopRoutes.PUT("/:shop_id", handler.UpdateShop)   // PUT /v1/shop/:shop_id
			shopRoutes.DELETE("/:shop_id", handler.DeleteShop) // DELETE /v1/shop/:shop_id
		}
	}
	return r
}
