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
			customers.POST("/register", handler.Register)
			customers.POST("/login", handler.Login)
			customers.POST("/logout", handler.Logout)

			protected := customers.Group("/")
			protected.Use(AuthMiddleware())
			{
				protected.GET("/me", handler.Me)
			}
		}
	}

	return r
}
