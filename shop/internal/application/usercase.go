package application

import "github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/domain"

// ShopUsecase defines the interface for shop-related business logic.
type ShopUsecase interface {
	CreateShop(name, ownerID, address string) (*domain.Shop, error)
	GetShopByID(id string) (*domain.Shop, error)
	GetShopsByOwnerID(ownerID string) ([]*domain.Shop, error)
	UpdateShop(id, name, address string, isActive bool) (*domain.Shop, error)
	DeleteShop(id string) error
}
