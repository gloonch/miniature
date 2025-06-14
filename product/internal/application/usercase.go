package application

import "miniature/product/internal/domain"

type Usecase interface {
	CreateProduct(shopIDStr, name, description string, price float64, sku string, stockQuantity int, creatingUserIDStr string) (*domain.Product, error)
	GetProductByID(id string) (*domain.Product, error)
	GetProductsByShopID(shopIDStr string /*, requestingUserIDStr string - for future auth */) ([]*domain.Product, error)
	UpdateProduct(productIDStr string, name *string, description *string, price *float64, sku *string, stockQuantity *int, isActive *bool, requestingUserIDStr string) (*domain.Product, error)
	DeleteProduct(productIDStr string, requestingUserIDStr string) error
}
