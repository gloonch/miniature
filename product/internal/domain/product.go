package domain

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID            uuid.UUID `json:"id"`
	ShopID        uuid.UUID `json:"shop_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"` // Consider using a specific decimal type for currency
	SKU           string    `json:"sku"`
	StockQuantity int       `json:"stock_quantity"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}
