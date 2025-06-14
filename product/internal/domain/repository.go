package domain

type Repository interface {
	Create(product *Product) error
	FindByID(id string) (*Product, error)
	FindByShopID(shopID string) ([]*Product, error)
	Update(product *Product) error
	Delete(id string) error
}

// ShopOwnershipCheckerRepository defines an interface for checking shop ownership.
// This is used by the product service to authorize actions on products based on shop ownership.
type ShopOwnershipCheckerRepository interface {
	IsShopOwner(userID, shopID string) (bool, error)
}
