package application

import "miniature/shop/internal/domain"

// Usecase defines the interface for shop-related business logic.
type Usecase interface {
	CreateShop(name, ownerID string) (*domain.Shop, error)
	GetShopByID(id string) (*domain.Shop, error)
	GetShopsByOwnerID(ownerID string) ([]*domain.Shop, error)
	UpdateShop(id, userIDFromTokenStr, name string, isActive bool) (*domain.Shop, error)
	DeleteShop(id, userIDFromTokenStr string) error
}
