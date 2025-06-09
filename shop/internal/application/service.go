package application

import (
	"errors" // Add import for errors
	"time"

	"github.com/google/uuid"
	"github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/domain"
)

type shopService struct {
	repo domain.ShopRepository
}

func NewShopService(repo domain.ShopRepository) ShopUsecase {
	return &shopService{repo: repo}
}

func (s *shopService) CreateShop(name, ownerIDStr, address string) (*domain.Shop, error) {
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		return nil, err // Or a more specific error
	}

	shop := &domain.Shop{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   ownerID,
		Address:   address,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.Create(shop)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

// Placeholders for other usecase methods
func (s *shopService) GetShopByID(id string) (*domain.Shop, error) {
	return s.repo.FindByID(id)
}

func (s *shopService) GetShopsByOwnerID(ownerID string) ([]*domain.Shop, error) {
	return s.repo.FindByOwnerID(ownerID)
}

func (s *shopService) UpdateShop(id, userIDFromTokenStr, name, address string, isActive bool) (*domain.Shop, error) {
	shop, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("database error while finding shop: " + err.Error())
	}
	if shop == nil {
		return nil, errors.New("shop not found") // Or a specific domain.ErrShopNotFound
	}

	// Ownership Check
	if shop.OwnerID.String() != userIDFromTokenStr {
		return nil, errors.New("user is not authorized to update this shop") // Or domain.ErrForbidden
	}

	// Update fields
	shop.Name = name
	shop.Address = address
	shop.IsActive = isActive
	shop.UpdatedAt = time.Now()

	err = s.repo.Update(shop)
	if err != nil {
		return nil, errors.New("database error while updating shop: " + err.Error())
	}
	return shop, nil
}

func (s *shopService) DeleteShop(id, userIDFromTokenStr string) error {
	shop, err := s.repo.FindByID(id) // Fetch shop to check ownership
	if err != nil {
		return errors.New("database error while finding shop: " + err.Error())
	}
	if shop == nil {
		return errors.New("shop not found") // Or domain.ErrShopNotFound
	}

	// Ownership Check
	if shop.OwnerID.String() != userIDFromTokenStr {
		return errors.New("user is not authorized to delete this shop") // Or domain.ErrForbidden
	}

	return s.repo.Delete(id)
}
