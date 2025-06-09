package interfaces

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *ProductHandler, authMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	v1 := r.Group("/v1")
	{
		shopProducts := v1.Group("/shops/:shop_id/products")
		if authMiddleware != nil {
			shopProducts.Use(authMiddleware)
		}
		if handler != nil {
			shopProducts.POST("", handler.CreateProduct)
			shopProducts.GET("", handler.GetShopProducts)
		}

		productRoutes := v1.Group("/products")
		if authMiddleware != nil {
			productRoutes.Use(authMiddleware)
		}
		if handler != nil {
			productRoutes.GET("/:product_id", handler.GetProduct)
			productRoutes.PUT("/:product_id", handler.UpdateProduct)
			productRoutes.DELETE("/:product_id", handler.DeleteProduct) // Added
		}
	}
	return r
}
