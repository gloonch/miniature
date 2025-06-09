package application

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/domain"
)

// --- MockShopRepository ---
type MockShopRepository struct {
	mock.Mock
}

func (m *MockShopRepository) Create(shop *domain.Shop) error {
	args := m.Called(shop)
	return args.Error(0)
}

func (m *MockShopRepository) FindByID(id string) (*domain.Shop, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Shop), args.Error(1)
}

func (m *MockShopRepository) FindByOwnerID(ownerID string) ([]*domain.Shop, error) {
	args := m.Called(ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Shop), args.Error(1)
}

func (m *MockShopRepository) Update(shop *domain.Shop) error {
	args := m.Called(shop)
	return args.Error(0)
}

func (m *MockShopRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- Tests for CreateShop ---
func TestCreateShop_Success(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)

	name := "Test Shop"
	ownerIDStr := uuid.New().String()
	address := "123 Test St"

	// Mock setup
	mockRepo.On("Create", mock.AnythingOfType("*domain.Shop")).Run(func(args mock.Arguments) {
		shopArg := args.Get(0).(*domain.Shop)
		assert.Equal(t, name, shopArg.Name)
		assert.Equal(t, ownerIDStr, shopArg.OwnerID.String())
		assert.Equal(t, address, shopArg.Address)
		assert.True(t, shopArg.IsActive)
		assert.NotNil(t, shopArg.ID)
		assert.WithinDuration(t, time.Now(), shopArg.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), shopArg.UpdatedAt, time.Second)
	}).Return(nil).Once()

	shop, err := service.CreateShop(name, ownerIDStr, address)

	assert.NoError(t, err)
	assert.NotNil(t, shop)
	assert.Equal(t, name, shop.Name)
	assert.Equal(t, ownerIDStr, shop.OwnerID.String())
	assert.Equal(t, address, shop.Address)
	assert.True(t, shop.IsActive)
	// ID, CreatedAt, UpdatedAt are checked by the mock.Run function
	mockRepo.AssertExpectations(t)
}

func TestCreateShop_RepositoryError(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	expectedErr := errors.New("repo create error")

	mockRepo.On("Create", mock.AnythingOfType("*domain.Shop")).Return(expectedErr).Once()

	shop, err := service.CreateShop("Test Shop", uuid.New().String(), "123 Test St")

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, shop)
	mockRepo.AssertExpectations(t)
}

// --- Tests for GetShopsByOwnerID ---
func TestGetShopsByOwnerID_Success(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	ownerID := uuid.New().String()
	expectedShops := []*domain.Shop{
		{ID: uuid.New(), Name: "Shop 1", OwnerID: uuid.MustParse(ownerID)},
		{ID: uuid.New(), Name: "Shop 2", OwnerID: uuid.MustParse(ownerID)},
	}

	mockRepo.On("FindByOwnerID", ownerID).Return(expectedShops, nil).Once()

	shops, err := service.GetShopsByOwnerID(ownerID)

	assert.NoError(t, err)
	assert.Equal(t, expectedShops, shops)
	mockRepo.AssertExpectations(t)
}

func TestGetShopsByOwnerID_Empty(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	ownerID := uuid.New().String()
	expectedShops := []*domain.Shop{} // Empty slice

	mockRepo.On("FindByOwnerID", ownerID).Return(expectedShops, nil).Once()

	shops, err := service.GetShopsByOwnerID(ownerID)

	assert.NoError(t, err)
	assert.Len(t, shops, 0)
	mockRepo.AssertExpectations(t)
}


func TestGetShopsByOwnerID_RepositoryError(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	ownerID := uuid.New().String()
	expectedErr := errors.New("repo find by owner error")

	mockRepo.On("FindByOwnerID", ownerID).Return(nil, expectedErr).Once()

	shops, err := service.GetShopsByOwnerID(ownerID)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, shops)
	mockRepo.AssertExpectations(t)
}


// --- Tests for UpdateShop ---
func TestUpdateShop_Success(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)

	shopID := uuid.New()
	ownerID := uuid.New()
	originalShop := &domain.Shop{
		ID:        shopID,
		Name:      "Original Name",
		OwnerID:   ownerID,
		Address:   "Old Address",
		IsActive:  true,
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	newName := "Updated Name"
	newAddress := "New Address"
	newIsActive := false

	mockRepo.On("FindByID", shopID.String()).Return(originalShop, nil).Once()
	mockRepo.On("Update", mock.MatchedBy(func(s *domain.Shop) bool {
		return s.ID == shopID && s.Name == newName && s.Address == newAddress && s.IsActive == newIsActive && s.OwnerID == ownerID
	})).Return(nil).Once()

	updatedShop, err := service.UpdateShop(shopID.String(), ownerID.String(), newName, newAddress, newIsActive)

	assert.NoError(t, err)
	assert.NotNil(t, updatedShop)
	assert.Equal(t, newName, updatedShop.Name)
	assert.Equal(t, newAddress, updatedShop.Address)
	assert.Equal(t, newIsActive, updatedShop.IsActive)
	assert.True(t, updatedShop.UpdatedAt.After(originalShop.UpdatedAt))
	mockRepo.AssertExpectations(t)
}

func TestUpdateShop_NotFound(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New().String()

	mockRepo.On("FindByID", shopID).Return(nil, nil).Once() // Not found

	shop, err := service.UpdateShop(shopID, uuid.New().String(), "N", "A", true)

	assert.Error(t, err)
	// Check for specific error message if service returns one, e.g., "shop not found"
	assert.Contains(t, err.Error(), "shop not found")
	assert.Nil(t, shop)
	mockRepo.AssertExpectations(t)
}

func TestUpdateShop_Forbidden(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)

	shopID := uuid.New()
	actualOwnerID := uuid.New()
	attackerOwnerID := uuid.New() // Different owner

	originalShop := &domain.Shop{ID: shopID, OwnerID: actualOwnerID, Name: "Test"}
	mockRepo.On("FindByID", shopID.String()).Return(originalShop, nil).Once()

	shop, err := service.UpdateShop(shopID.String(), attackerOwnerID.String(), "N", "A", true)

	assert.Error(t, err)
	// Check for specific error message, e.g., "user is not authorized"
	assert.Contains(t, err.Error(), "user is not authorized")
	assert.Nil(t, shop)
	mockRepo.AssertExpectations(t) // Verifies FindByID was called
	mockRepo.AssertNotCalled(t, "Update", mock.Anything) // Crucially, Update should not be called
}

func TestUpdateShop_FindByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New().String()
	ownerID := uuid.New().String()
	expectedErr := errors.New("db find error")

	mockRepo.On("FindByID", shopID).Return(nil, expectedErr).Once()

	shop, err := service.UpdateShop(shopID, ownerID, "N", "A", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error while finding shop")
	assert.Nil(t, shop)
	mockRepo.AssertExpectations(t)
}


func TestUpdateShop_Update_RepositoryError(t *testing.T) {
    mockRepo := new(MockShopRepository)
    service := NewShopService(mockRepo)

    shopID := uuid.New()
    ownerID := uuid.New()
    originalShop := &domain.Shop{ID: shopID, Name: "Original", OwnerID: ownerID}
    expectedErr := errors.New("repo update error")

    mockRepo.On("FindByID", shopID.String()).Return(originalShop, nil).Once()
    mockRepo.On("Update", mock.AnythingOfType("*domain.Shop")).Return(expectedErr).Once()

    shop, err := service.UpdateShop(shopID.String(), ownerID.String(), "New Name", "New Address", true)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database error while updating shop") // Based on current service error
    assert.Nil(t, shop)
    mockRepo.AssertExpectations(t)
}


// --- Tests for DeleteShop ---
func TestDeleteShop_Success(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New()
	ownerID := uuid.New()
	existingShop := &domain.Shop{ID: shopID, OwnerID: ownerID}

	mockRepo.On("FindByID", shopID.String()).Return(existingShop, nil).Once()
	mockRepo.On("Delete", shopID.String()).Return(nil).Once()

	err := service.DeleteShop(shopID.String(), ownerID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteShop_NotFound(t *testing.T) {
    mockRepo := new(MockShopRepository)
    service := NewShopService(mockRepo)
    shopID := uuid.New().String()
    ownerID := uuid.New().String()

    mockRepo.On("FindByID", shopID).Return(nil, nil).Once() // Simulate shop not found

    err := service.DeleteShop(shopID, ownerID)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "shop not found")
    mockRepo.AssertExpectations(t)
    mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
}


func TestDeleteShop_Forbidden(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New()
	actualOwnerID := uuid.New()
	attackerOwnerID := uuid.New()
	existingShop := &domain.Shop{ID: shopID, OwnerID: actualOwnerID}

	mockRepo.On("FindByID", shopID.String()).Return(existingShop, nil).Once()

	err := service.DeleteShop(shopID.String(), attackerOwnerID.String())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user is not authorized")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
}

func TestDeleteShop_FindByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New().String()
	ownerID := uuid.New().String()
	expectedErr := errors.New("db find error")

	mockRepo.On("FindByID", shopID).Return(nil, expectedErr).Once()

	err := service.DeleteShop(shopID, ownerID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error while finding shop")
	mockRepo.AssertExpectations(t)
}


func TestDeleteShop_Delete_RepositoryError(t *testing.T) {
    mockRepo := new(MockShopRepository)
    service := NewShopService(mockRepo)
    shopID := uuid.New()
    ownerID := uuid.New()
    existingShop := &domain.Shop{ID: shopID, OwnerID: ownerID}
    expectedErr := errors.New("repo delete error")

    mockRepo.On("FindByID", shopID.String()).Return(existingShop, nil).Once()
    mockRepo.On("Delete", shopID.String()).Return(expectedErr).Once()

    err := service.DeleteShop(shopID.String(), ownerID.String())

    assert.Error(t, err)
    assert.EqualError(t, err, expectedErr.Error()) // Service directly returns repo error here
    mockRepo.AssertExpectations(t)
}

// --- Tests for GetShopByID ---
func TestGetShopByID_Success(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New().String()
	expectedShop := &domain.Shop{ID: uuid.MustParse(shopID), Name: "Found Shop"}

	mockRepo.On("FindByID", shopID).Return(expectedShop, nil).Once()

	shop, err := service.GetShopByID(shopID)

	assert.NoError(t, err)
	assert.Equal(t, expectedShop, shop)
	mockRepo.AssertExpectations(t)
}

func TestGetShopByID_NotFound(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New().String()

	mockRepo.On("FindByID", shopID).Return(nil, nil).Once() // Repository returns nil, nil for not found

	shop, err := service.GetShopByID(shopID)

	assert.NoError(t, err) // The service currently returns nil, nil for not found
	assert.Nil(t, shop)
	mockRepo.AssertExpectations(t)
}

func TestGetShopByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockShopRepository)
	service := NewShopService(mockRepo)
	shopID := uuid.New().String()
	expectedErr := errors.New("repo find error")

	mockRepo.On("FindByID", shopID).Return(nil, expectedErr).Once()

	shop, err := service.GetShopByID(shopID)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())
	assert.Nil(t, shop)
	mockRepo.AssertExpectations(t)
}
