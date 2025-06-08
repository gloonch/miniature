package interfaces

import (
	"github.com/gin-gonic/gin"
	"miniature/customer/internal/application"
	"net/http"
)

type CustomerHandler struct {
	usecase application.CustomerUsecase
}

func NewCustomerHandler(u application.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: u}
}

func (h *CustomerHandler) Register(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})

		return
	}

	customer, cErr := h.usecase.RegisterCustomer(req.Name, req.Phone, req.Role)
	if cErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": cErr.Error()})

		return
	}

	c.JSON(http.StatusCreated, customer)
}
