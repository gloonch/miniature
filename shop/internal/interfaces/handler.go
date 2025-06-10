package interfaces

import (
	"database/sql" // Added for sql.ErrNoRows
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/application"
	// "github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/domain" // For ShopResponse, if it's different from domain.Shop
)

type ShopHandler struct {
	usecase application.ShopUsecase
}

func NewShopHandler(u application.ShopUsecase) *ShopHandler {
	return &ShopHandler{usecase: u}
}

func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req CreateShopRequest
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
	roleStr, ok := roleRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role is not of type string"})
		return
	}

	// Authorization: Only "SELLER" role can create shops for now
	// TODO: Make "SELLER" a constant or configurable
	if roleStr != "SELLER" {
		c.JSON(http.StatusForbidden, gin.H{"error": "user does not have permission to create a shop"})
		return
	}

	shop, err := h.usecase.CreateShop(req.Name, userIDStr, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create shop: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shop) // Assuming domain.Shop can be directly marshalled
}

// Placeholder handlers
func (h *ShopHandler) GetShop(c *gin.Context) {
	shopID := c.Param("shop_id")

	shop, err := h.usecase.GetShopByID(shopID)
	if err != nil {
		// Assuming service layer might return specific errors like not found
		// For now, a general server error if not nil, or check if shop is nil for not found.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve shop: " + err.Error()})
		return
	}
	if shop == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "shop not found"})
		return
	}
	c.JSON(http.StatusOK, shop)
}

func (h *ShopHandler) GetUserShops(c *gin.Context) {
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

	shops, err := h.usecase.GetShopsByOwnerID(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve user shops: " + err.Error()})
		return
	}
	// Always return a list, even if empty
	c.JSON(http.StatusOK, shops)
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {
	shopID := c.Param("shop_id")

	var req UpdateShopRequest
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
	roleStr, ok := roleRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role is not of type string"})
		return
	}

	if roleStr != "SELLER" { // Assuming "SELLER" role is required to update
		c.JSON(http.StatusForbidden, gin.H{"error": "user does not have permission to update shops"})
		return
	}

	updatedShop, err := h.usecase.UpdateShop(shopID, userIDStr, req.Name, req.Address, req.IsActive)
	if err != nil {
		// Basic error handling, can be more granular
		if err.Error() == "shop not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "user is not authorized to update this shop" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update shop: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedShop)
}

func (h *ShopHandler) DeleteShop(c *gin.Context) {
	shopID := c.Param("shop_id")

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
	roleStr, ok := roleRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role is not of type string"})
		return
	}

	if roleStr != "SELLER" { // Assuming "SELLER" role is required to delete
		c.JSON(http.StatusForbidden, gin.H{"error": "user does not have permission to delete shops"})
		return
	}

	err := h.usecase.DeleteShop(shopID, userIDStr)
	if err != nil {
		if err.Error() == "shop not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "user is not authorized to delete this shop" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		// Check for sql.ErrNoRows from repository delete if it means not found
		if err.Error() == sql.ErrNoRows.Error() { // Comparing error strings, better to use errors.Is in Go 1.13+
			c.JSON(http.StatusNotFound, gin.H{"error": "shop not found to delete"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete shop: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
