package interfaces

import (
	"time"

	"github.com/google/uuid"
)

// CreateShopRequest represents the request payload for creating a new shop.
type CreateShopRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateShopRequest represents the request payload for updating an existing shop.
type UpdateShopRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// ShopResponse represents the response payload for a shop.
type ShopResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
