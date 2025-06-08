package interfaces

import (
	"github.com/gin-gonic/gin"
	"miniature/customer/internal/application"
	"miniature/pkg/token"
	"net/http"
)

type CustomerHandler struct {
	usecase application.CustomerUsecase
}

func NewCustomerHandler(u application.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: u}
}

func (h *CustomerHandler) Register(c *gin.Context) {
	var req CreateCustomerRequest

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

func (h *CustomerHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	customer, err := h.usecase.GetCustomerByPhone(req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	token, err := token.GenerateToken(customer.ID.String(), customer.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *CustomerHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout successful. Just discard the token on client side."})
}

func (h *CustomerHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	customer, err := h.usecase.GetCustomerByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}
