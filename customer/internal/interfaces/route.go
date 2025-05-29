package interfaces

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler CustomerHandler) *gin.Engine {
	r := gin.Default()

	// Grouped routes
	v1 := r.Group("/v1")
	{
		customers := v1.Group("/customer")
		{
			customers.POST("/", handler.Register)
			// Add more routes...
		}
	}

	return r
}
