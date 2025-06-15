package interfaces

import (
	"database/sql" // For sql.ErrNoRows check from service
	"miniature/product/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase application.Usecase
}

func NewHandler(u application.Usecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) CreateProduct(c *gin.Context) {
	shopIDStr := c.Param("shop_id")
	// Potentially validate shopIDStr format here if not done by a path regex

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id is not of type string"})
		return
	}

	roleRaw, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in context"})
		return
	}

	_, ok = roleRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role is not of type string"})
		return
	}

	product, err := h.usecase.CreateProduct(shopIDStr, req.Name, req.Description, req.Price, req.SKU, req.StockQuantity, userIDStr)
	if err != nil {
		// Check for specific errors from usecase
		if err.Error() == "user not authorized to add products to this shop" ||
			err.Error() == "could not verify shop ownership" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		// if strings.Contains(err.Error(), "already exists") { // For SKU conflict
		// 	c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		// 	return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create product: " + err.Error()})
		return
	}

	// Convert domain.Product to ProductResponse if they are different, or use domain.Product directly
	c.JSON(http.StatusCreated, product)
}

func (h *Handler) GetProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")

	product, err := h.usecase.GetProductByID(productIDStr)
	if err != nil {
		// Assuming service layer might return specific errors like not found
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve product: " + err.Error()})
		return
	}
	if product == nil { // Standard check for not found if service returns (nil, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product) // Using domain.Product as response for now
}

func (h *Handler) GetShopProducts(c *gin.Context) {
	shopIDStr := c.Param("shop_id")

	// userIDRaw, _ := c.Get("user_id") // For future authorization if needed
	// userIDStr, _ := userIDRaw.(string)

	products, err := h.usecase.GetProductsByShopID(shopIDStr /*, userIDStr */)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve shop products: " + err.Error()})
		return
	}
	// Always return a list, even if empty
	c.JSON(http.StatusOK, products) // Using slice of domain.Product as response
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists { // Should be caught by AuthMiddleware, but good practice
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userIDStr, _ := userIDRaw.(string)
	// No role check here, service layer should check ownership via shop

	updatedProduct, err := h.usecase.UpdateProduct(
		productIDStr,
		req.Name,
		req.Description,
		req.Price,
		req.SKU,
		req.StockQuantity,
		req.IsActive,
		userIDStr,
	)

	if err != nil {
		// Example error handling, can be more granular based on errors from usecase
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "user not authorized to update this product" ||
			err.Error() == "could not verify shop ownership for product update" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		// if strings.Contains(err.Error(), "already exists") { // For SKU conflict
		//  c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		//  return
		// }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	userIDStr, _ := userIDRaw.(string)

	err := h.usecase.DeleteProduct(productIDStr, userIDStr)

	if err != nil {
		if err == sql.ErrNoRows { // Check for not found error from service
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err.Error() == "user not authorized to delete this product" ||
			err.Error() == "could not verify shop ownership for product deletion" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
