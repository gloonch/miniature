package interfaces

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	v1 := r.Group("/v1")
	{
		shopProducts := v1.Group("/shops/:shop_id/products")
		shopProducts.Use(AuthMiddleware())
		{
			shopProducts.POST("", handler.CreateProduct)
			shopProducts.GET("", handler.GetShopProducts)
		}

		productRoutes := v1.Group("/products")
		productRoutes.Use(AuthMiddleware())
		{
			productRoutes.GET("/:product_id", handler.GetProduct)
			productRoutes.PUT("/:product_id", handler.UpdateProduct)
			productRoutes.DELETE("/:product_id", handler.DeleteProduct)
		}
	}
	return r
}
