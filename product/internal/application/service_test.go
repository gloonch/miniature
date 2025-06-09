package application

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/segment-sources/sources-backend-takehome-assignment/product/internal/domain"
)

// --- MockProductRepository ---
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductRepository) FindByID(id string) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}
func (m *MockProductRepository) FindByShopID(shopID string) ([]*domain.Product, error) {
	args := m.Called(shopID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}
func (m *MockProductRepository) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- MockShopOwnershipCheckerRepository ---
type MockShopOwnershipCheckerRepository struct {
	mock.Mock
}

func (m *MockShopOwnershipCheckerRepository) IsShopOwner(userID, shopID string) (bool, error) {
	args := m.Called(userID, shopID)
	return args.Bool(0), args.Error(1)
}

// --- Tests for CreateProduct ---
func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)

	shopIDStr := uuid.New().String()
	userIDStr := uuid.New().String()
	reqName, reqDesc, reqPrice, reqSKU, reqStock := "Test Prod", "Desc", 10.0, "SKU123", 5

	mockShopChecker.On("IsShopOwner", userIDStr, shopIDStr).Return(true, nil).Once()
	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Run(func(args mock.Arguments) {
		p := args.Get(0).(*domain.Product)
		assert.Equal(t, reqName, p.Name)
		assert.Equal(t, shopIDStr, p.ShopID.String())
		assert.Equal(t, reqDesc, p.Description)
		assert.Equal(t, reqPrice, p.Price)
		assert.Equal(t, reqSKU, p.SKU)
		assert.Equal(t, reqStock, p.StockQuantity)
		assert.True(t, p.IsActive)
		assert.NotNil(t, p.ID)
		assert.WithinDuration(t, time.Now(), p.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), p.UpdatedAt, time.Second)
	}).Return(nil).Once()

	product, err := service.CreateProduct(shopIDStr, reqName, reqDesc, reqPrice, reqSKU, reqStock, userIDStr)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, reqName, product.Name)
	mockRepo.AssertExpectations(t)
	mockShopChecker.AssertExpectations(t)
}

func TestCreateProduct_ShopOwnershipCheckFails_NotOwner(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopIDStr, userIDStr := uuid.New().String(), uuid.New().String()

	mockShopChecker.On("IsShopOwner", userIDStr, shopIDStr).Return(false, nil).Once()

	_, err := service.CreateProduct(shopIDStr, "N", "D", 1.0, "S", 1, userIDStr)
	assert.Error(t, err)
	assert.Equal(t, "user not authorized to add products to this shop", err.Error())
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestCreateProduct_ShopOwnershipCheckError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopIDStr, userIDStr := uuid.New().String(), uuid.New().String()
	expectedErr := errors.New("db error checking ownership")

	mockShopChecker.On("IsShopOwner", userIDStr, shopIDStr).Return(false, expectedErr).Once()

	_, err := service.CreateProduct(shopIDStr, "N", "D", 1.0, "S", 1, userIDStr)
	assert.Error(t, err)
	assert.Equal(t, "could not verify shop ownership", err.Error())
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestCreateProduct_InvalidPrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopIDStr, userIDStr := uuid.New().String(), uuid.New().String()

	// No mocks expected to be called for this validation error
	_, err := service.CreateProduct(shopIDStr, "N", "D", -1.0, "S", 1, userIDStr)
	assert.Error(t, err)
	assert.Equal(t, "price cannot be negative", err.Error())
}

func TestCreateProduct_InvalidStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopIDStr, userIDStr := uuid.New().String(), uuid.New().String()

	// Price is valid, but stock is not
	_, err := service.CreateProduct(shopIDStr, "N", "D", 1.0, "S", -1, userIDStr)
	assert.Error(t, err)
	assert.Equal(t, "stock quantity cannot be negative", err.Error())
}

func TestCreateProduct_RepoCreateError_SKUConflict(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopIDStr, userIDStr := uuid.New().String(), uuid.New().String()
	// This error message is a placeholder; actual message depends on DB driver and error type.
	// The service currently doesn't parse this specific error, so it's treated as a generic repo error.
	expectedErr := errors.New("pq: duplicate key value violates unique constraint uq_shop_sku")

	mockShopChecker.On("IsShopOwner", userIDStr, shopIDStr).Return(true, nil).Once()
	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(expectedErr).Once()

	_, err := service.CreateProduct(shopIDStr, "N", "D", 1.0, "S", 1, userIDStr)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err) // Service propagates the error as is
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_RepoCreateError_Generic(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopIDStr, userIDStr := uuid.New().String(), uuid.New().String()
	expectedErr := errors.New("generic db error")

	mockShopChecker.On("IsShopOwner", userIDStr, shopIDStr).Return(true, nil).Once()
	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(expectedErr).Once()

	_, err := service.CreateProduct(shopIDStr, "N", "D", 1.0, "S", 1, userIDStr)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}


// --- Tests for GetProductByID ---
func TestGetProductByID_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productID := uuid.New().String()
	expectedProduct := &domain.Product{ID: uuid.MustParse(productID), Name: "Found Product"}

	mockRepo.On("FindByID", productID).Return(expectedProduct, nil).Once()
	product, err := service.GetProductByID(productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productID := uuid.New().String()

	mockRepo.On("FindByID", productID).Return(nil, nil).Once() // Repo returns nil, nil for not found
	product, err := service.GetProductByID(productID)

	assert.NoError(t, err) // Service returns nil, nil as well
	assert.Nil(t, product)
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_RepoError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productID := uuid.New().String()
	expectedErr := errors.New("repo findbyid error")

	mockRepo.On("FindByID", productID).Return(nil, expectedErr).Once()
	product, err := service.GetProductByID(productID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, product)
	mockRepo.AssertExpectations(t)
}

// --- Tests for GetProductsByShopID ---
func TestGetProductsByShopID_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopID := uuid.New().String()
	expectedProducts := []*domain.Product{
		{ID: uuid.New(), Name: "Prod 1", ShopID: uuid.MustParse(shopID)},
		{ID: uuid.New(), Name: "Prod 2", ShopID: uuid.MustParse(shopID)},
	}
	mockRepo.On("FindByShopID", shopID).Return(expectedProducts, nil).Once()
	products, err := service.GetProductsByShopID(shopID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByShopID_Empty(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopID := uuid.New().String()
	expectedProducts := []*domain.Product{} // Empty slice

	mockRepo.On("FindByShopID", shopID).Return(expectedProducts, nil).Once()
	products, err := service.GetProductsByShopID(shopID)

	assert.NoError(t, err)
	assert.Len(t, products, 0)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByShopID_RepoError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	shopID := uuid.New().String()
	expectedErr := errors.New("repo findbyshopid error")

	mockRepo.On("FindByShopID", shopID).Return(nil, expectedErr).Once()
	products, err := service.GetProductsByShopID(shopID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, products)
	mockRepo.AssertExpectations(t)
}

// --- Helper function for pointer values ---
func ptr[T any](v T) *T {
	return &v
}

// --- Tests for UpdateProduct ---
func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)

	productID := uuid.New()
	shopID := uuid.New()
	userID := uuid.New().String()

	existingProduct := &domain.Product{
		ID: productID, ShopID: shopID, Name: "Old Name", Price: 10.0, SKU: "OLD_SKU", StockQuantity: 10, IsActive: true, UpdatedAt: time.Now().Add(-time.Hour),
	}
	mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
	mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(true, nil).Once()
	mockRepo.On("Update", mock.MatchedBy(func(p *domain.Product) bool {
		// Check fields that are expected to be updated
		return p.ID == productID &&
			   p.Name == "New Name" && // Updated
			   p.Price == 20.0 && // Updated
			   p.SKU == "OLD_SKU" && // Unchanged
			   p.StockQuantity == 5 && // Updated
			   p.IsActive == false // Updated
	})).Return(nil).Once()

	updatedProduct, err := service.UpdateProduct(productID.String(),
		ptr("New Name"),                  // name
		nil,                             // description (not updated)
		ptr(20.0),                       // price
		nil,                             // sku (not updated)
		ptr(5),                          // stockQuantity
		ptr(false),                      // isActive
		userID)

	assert.NoError(t, err)
	assert.NotNil(t, updatedProduct)
	assert.Equal(t, "New Name", updatedProduct.Name)
	assert.Equal(t, 20.0, updatedProduct.Price)
	assert.Equal(t, 5, updatedProduct.StockQuantity)
	assert.False(t, updatedProduct.IsActive)
	assert.True(t, updatedProduct.UpdatedAt.After(existingProduct.UpdatedAt))
	mockRepo.AssertExpectations(t)
	mockShopChecker.AssertExpectations(t)
}

func TestUpdateProduct_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productIDStr := uuid.New().String()

	mockRepo.On("FindByID", productIDStr).Return(nil, nil).Once() // Product not found

	_, err := service.UpdateProduct(productIDStr, ptr("N"), nil, nil, nil, nil, nil, "userID")
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	mockShopChecker.AssertNotCalled(t, "IsShopOwner", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestUpdateProduct_OwnershipCheckFails_NotOwner(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
	existingProduct := &domain.Product{ID: productID, ShopID: shopID}

	mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
	mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(false, nil).Once()

	_, err := service.UpdateProduct(productID.String(), ptr("N"), nil,nil,nil,nil,nil, userID)
	assert.Error(t, err)
	assert.Equal(t, "user not authorized to update this product", err.Error())
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_OwnershipCheckError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
	existingProduct := &domain.Product{ID: productID, ShopID: shopID}
	expectedErr := errors.New("db error checking ownership")

	mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
	mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(false, expectedErr).Once()

	_, err := service.UpdateProduct(productID.String(), ptr("N"), nil,nil,nil,nil,nil, userID)
	assert.Error(t, err)
	assert.Equal(t, "could not verify shop ownership for product update", err.Error())
	mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_InvalidPriceUpdate(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockShopChecker := new(MockShopOwnershipCheckerRepository)
    service := NewProductService(mockRepo, mockShopChecker)
    productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
    existingProduct := &domain.Product{ID: productID, ShopID: shopID}

    mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
    mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(true, nil).Once()

    _, err := service.UpdateProduct(productID.String(), nil, nil, ptr(-5.0), nil, nil, nil, userID)
    assert.Error(t, err)
    assert.Equal(t, "price cannot be negative", err.Error())
    mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_InvalidStockUpdate(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockShopChecker := new(MockShopOwnershipCheckerRepository)
    service := NewProductService(mockRepo, mockShopChecker)
    productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
    existingProduct := &domain.Product{ID: productID, ShopID: shopID}

    mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
    mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(true, nil).Once()

    _, err := service.UpdateProduct(productID.String(), nil, nil, nil, nil, ptr(-10), nil, userID)
    assert.Error(t, err)
    assert.Equal(t, "stock quantity cannot be negative", err.Error())
    mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_RepoUpdateError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
	existingProduct := &domain.Product{ID: productID, ShopID: shopID}
	expectedErr := errors.New("db update error")

	mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
	mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(true, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*domain.Product")).Return(expectedErr).Once()

	_, err := service.UpdateProduct(productID.String(), ptr("N"), nil,nil,nil,nil,nil, userID)
	assert.Error(t, err)
	assert.Equal(t, "database error while updating product: "+expectedErr.Error(), err.Error())
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}


// --- Tests for DeleteProduct ---
func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
	existingProduct := &domain.Product{ID: productID, ShopID: shopID}

	mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
	mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(true, nil).Once()
	mockRepo.On("Delete", productID.String()).Return(nil).Once()

	err := service.DeleteProduct(productID.String(), userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockShopChecker.AssertExpectations(t)
}

func TestDeleteProduct_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockShopChecker := new(MockShopOwnershipCheckerRepository)
	service := NewProductService(mockRepo, mockShopChecker)
	productIDStr := uuid.New().String()

	mockRepo.On("FindByID", productIDStr).Return(nil, sql.ErrNoRows).Once()

	err := service.DeleteProduct(productIDStr, "userID")
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err) // Service propagates sql.ErrNoRows
	mockShopChecker.AssertNotCalled(t, "IsShopOwner", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
}

func TestDeleteProduct_OwnershipCheckFails_NotOwner(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockShopChecker := new(MockShopOwnershipCheckerRepository)
    service := NewProductService(mockRepo, mockShopChecker)
    productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
    existingProduct := &domain.Product{ID: productID, ShopID: shopID}

    mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
    mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(false, nil).Once()

    err := service.DeleteProduct(productID.String(), userID)
    assert.Error(t, err)
    assert.Equal(t, "user not authorized to delete this product", err.Error())
    mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_OwnershipCheckError(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockShopChecker := new(MockShopOwnershipCheckerRepository)
    service := NewProductService(mockRepo, mockShopChecker)
    productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
    existingProduct := &domain.Product{ID: productID, ShopID: shopID}
	expectedErr := errors.New("db error checking ownership")

    mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
    mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(false, expectedErr).Once()

    err := service.DeleteProduct(productID.String(), userID)
    assert.Error(t, err)
    assert.Equal(t, "could not verify shop ownership for product deletion", err.Error())
    mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_RepoDeleteError(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockShopChecker := new(MockShopOwnershipCheckerRepository)
    service := NewProductService(mockRepo, mockShopChecker)
    productID := uuid.New(); shopID := uuid.New(); userID := uuid.New().String()
    existingProduct := &domain.Product{ID: productID, ShopID: shopID}
	expectedErr := errors.New("db delete error")

    mockRepo.On("FindByID", productID.String()).Return(existingProduct, nil).Once()
    mockShopChecker.On("IsShopOwner", userID, shopID.String()).Return(true, nil).Once()
	mockRepo.On("Delete", productID.String()).Return(expectedErr).Once()

    err := service.DeleteProduct(productID.String(), userID)
    assert.Error(t, err)
    assert.Equal(t, expectedErr, err.Error())
	mockShopChecker.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_FindByIDError(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockShopChecker := new(MockShopOwnershipCheckerRepository)
    service := NewProductService(mockRepo, mockShopChecker)
    productIDStr := uuid.New().String()
	expectedErr := errors.New("other db error")


    mockRepo.On("FindByID", productIDStr).Return(nil, expectedErr).Once()

    err := service.DeleteProduct(productIDStr, "userID")
    assert.Error(t, err)
	assert.Equal(t, "database error while finding product: "+expectedErr.Error(), err.Error())
    mockShopChecker.AssertNotCalled(t, "IsShopOwner", mock.Anything, mock.Anything)
    mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
	mockRepo.AssertExpectations(t)
}
