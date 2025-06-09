package interfaces

import (
	"time"
	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" binding:"required,gte=0"`
	SKU           string  `json:"sku"`
	StockQuantity int     `json:"stock_quantity" binding:"gte=0"`
}

// ProductResponse can be the domain.Product or a specific DTO
// For now, using domain.Product directly is fine for responses.
type ProductResponse struct {
	ID            uuid.UUID `json:"id"`
	ShopID        uuid.UUID `json:"shop_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	SKU           string    `json:"sku"`
	StockQuantity int       `json:"stock_quantity"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdateProductRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	Price         *float64 `json:"price" binding:"omitempty,gte=0"` // omitempty allows nil, gte=0 applies if not nil
	SKU           *string  `json:"sku"`
	StockQuantity *int     `json:"stock_quantity" binding:"omitempty,gte=0"`
	IsActive      *bool    `json:"is_active"`
}
