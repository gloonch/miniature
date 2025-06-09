package domain

// ShopRepository defines the interface for interacting with shop data.
type ShopRepository interface {
	Create(shop *Shop) error
	FindByID(id string) (*Shop, error)
	FindByOwnerID(ownerID string) ([]*Shop, error)
	Update(shop *Shop) error
	Delete(id string) error
}
