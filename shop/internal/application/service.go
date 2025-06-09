package application

import (
	"errors"
	"miniature/shop/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.Repository
}

func NewShopService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateShop(name, ownerIDStr string) (*domain.Shop, error) {
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		return nil, err // Or a more specific error
	}

	shop := &domain.Shop{
		ID:        uuid.New(),
		Name:      name,
		OwnerID:   ownerID,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	err = s.repo.Create(shop)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

// Placeholders for other usecase methods
func (s *Service) GetShopByID(id string) (*domain.Shop, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetShopsByOwnerID(ownerID string) ([]*domain.Shop, error) {
	return s.repo.FindByOwnerID(ownerID)
}

func (s *Service) UpdateShop(id, userIDFromTokenStr, name string, isActive bool) (*domain.Shop, error) {
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
	shop.IsActive = isActive

	err = s.repo.Update(shop)
	if err != nil {
		return nil, errors.New("database error while updating shop: " + err.Error())
	}
	return shop, nil
}

func (s *Service) DeleteShop(id, userIDFromTokenStr string) error {
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
