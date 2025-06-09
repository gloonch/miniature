package interfaces

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/application"
	"github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MockShopUsecase ---
type MockShopUsecase struct {
	mock.Mock
}

func (m *MockShopUsecase) CreateShop(name, ownerID, address string) (*domain.Shop, error) {
	args := m.Called(name, ownerID, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Shop), args.Error(1)
}

func (m *MockShopUsecase) GetShopByID(id string) (*domain.Shop, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Shop), args.Error(1)
}

func (m *MockShopUsecase) GetShopsByOwnerID(ownerID string) ([]*domain.Shop, error) {
	args := m.Called(ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Shop), args.Error(1)
}

func (m *MockShopUsecase) UpdateShop(id, userIDFromToken, name, address string, isActive bool) (*domain.Shop, error) {
	args := m.Called(id, userIDFromToken, name, address, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Shop), args.Error(1)
}

func (m *MockShopUsecase) DeleteShop(id, userIDFromToken string) error {
	args := m.Called(id, userIDFromToken)
	return args.Error(0)
}

// Helper to perform requests
func performHandlerTestRequest(t *testing.T, r *gin.Engine, method, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, body)
	assert.NoError(t, err)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

// Middleware factory for setting test context values
func testContextMiddleware(userID, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	}
}

// --- Tests for CreateShop Handler ---
func TestCreateShop_Handler_Success(t *testing.T) {
	mockUsecase := new(MockShopUsecase)
	handler := NewShopHandler(mockUsecase)
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userID := uuid.New().String()
	r.Use(testContextMiddleware(userID, "SELLER"))
	r.POST("/v1/shop", handler.CreateShop)

	expectedShop := &domain.Shop{ID: uuid.New(), Name: "Test Shop", OwnerID: uuid.MustParse(userID), Address: "123 St"}
	createReq := CreateShopRequest{Name: "Test Shop", Address: "123 St"}
	mockUsecase.On("CreateShop", createReq.Name, userID, createReq.Address).Return(expectedShop, nil).Once()

	jsonBody, _ := json.Marshal(createReq)
	rr := performHandlerTestRequest(t, r, http.MethodPost, "/v1/shop", bytes.NewBuffer(jsonBody), map[string]string{"Content-Type": "application/json"})

	assert.Equal(t, http.StatusCreated, rr.Code)
	var returnedShop domain.Shop
	err := json.Unmarshal(rr.Body.Bytes(), &returnedShop)
	assert.NoError(t, err)
	assert.Equal(t, expectedShop.Name, returnedShop.Name)
	assert.Equal(t, expectedShop.Address, returnedShop.Address)
	mockUsecase.AssertExpectations(t)
}

func TestCreateShop_Handler_InvalidInput(t *testing.T) {
	mockUsecase := new(MockShopUsecase)
	handler := NewShopHandler(mockUsecase)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(testContextMiddleware(uuid.New().String(), "SELLER"))
	r.POST("/v1/shop", handler.CreateShop)

	// Corrected the map literal here
	rr := performHandlerTestRequest(t, r, http.MethodPost, "/v1/shop", bytes.NewBufferString(`{"name":`), map[string]string{"Content-Type": "application/json"})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertNotCalled(t, "CreateShop", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateShop_Handler_Forbidden(t *testing.T) {
	mockUsecase := new(MockShopUsecase)
	handler := NewShopHandler(mockUsecase)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(testContextMiddleware(uuid.New().String(), "CUSTOMER")) // Non-SELLER role
	r.POST("/v1/shop", handler.CreateShop)

	createReq := CreateShopRequest{Name: "Test Shop", Address: "123 St"}
	jsonBody, _ := json.Marshal(createReq)
	rr := performHandlerTestRequest(t, r, http.MethodPost, "/v1/shop", bytes.NewBuffer(jsonBody), map[string]string{"Content-Type": "application/json"})

	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockUsecase.AssertNotCalled(t, "CreateShop", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateShop_Handler_UsecaseError(t *testing.T) {
	mockUsecase := new(MockShopUsecase)
	handler := NewShopHandler(mockUsecase)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	userID := uuid.New().String()
	r.Use(testContextMiddleware(userID, "SELLER"))
	r.POST("/v1/shop", handler.CreateShop)

	createReq := CreateShopRequest{Name: "Test Shop", Address: "123 St"}
	mockUsecase.On("CreateShop", createReq.Name, userID, createReq.Address).Return(nil, errors.New("usecase failed")).Once()

	jsonBody, _ := json.Marshal(createReq)
	rr := performHandlerTestRequest(t, r, http.MethodPost, "/v1/shop", bytes.NewBuffer(jsonBody), map[string]string{"Content-Type": "application/json"})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockUsecase.AssertExpectations(t)
}

// --- Tests for GetShop Handler ---
func TestGetShop_Handler_Success(t *testing.T) {
	mockUsecase := new(MockShopUsecase)
	handler := NewShopHandler(mockUsecase)
	gin.SetMode(gin.TestMode)
	r := gin.New()

	shopID := uuid.New()
	// Middleware for user_id/role context. Path param shop_id is handled by Gin router.
	r.Use(testContextMiddleware(uuid.New().String(), "CUSTOMER"))
	r.GET("/v1/shop/:shop_id", handler.GetShop)

	expectedShop := &domain.Shop{ID: shopID, Name: "Found Shop"}
	mockUsecase.On("GetShopByID", shopID.String()).Return(expectedShop, nil).Once()

	// Path for request includes the actual shop_id for routing
	rr := performHandlerTestRequest(t, r, http.MethodGet, "/v1/shop/"+shopID.String(), nil, nil)

	assert.Equal(t, http.StatusOK, rr.Code)
	var returnedShop domain.Shop
	err := json.Unmarshal(rr.Body.Bytes(), &returnedShop)
	assert.NoError(t, err)
	assert.Equal(t, expectedShop.Name, returnedShop.Name)
	mockUsecase.AssertExpectations(t)
}

func TestGetShop_Handler_NotFound(t *testing.T) {
	mockUsecase := new(MockShopUsecase)
	handler := NewShopHandler(mockUsecase)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	shopIDstr := uuid.New().String()

	r.Use(testContextMiddleware(uuid.New().String(), "CUSTOMER"))
	r.GET("/v1/shop/:shop_id", handler.GetShop)

	// Usecase returns nil, nil to indicate not found (as per current service behavior)
	mockUsecase.On("GetShopByID", shopIDstr).Return(nil, nil).Once()

	rr := performHandlerTestRequest(t, r, http.MethodGet, "/v1/shop/"+shopIDstr, nil, nil)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockUsecase.AssertExpectations(t)
}

func TestGetShop_Handler_UsecaseError(t *testing.T) {
	mockUsecase := new(MockShopUsecase)
	handler := NewShopHandler(mockUsecase)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	shopIDstr := uuid.New().String()

	r.Use(testContextMiddleware(uuid.New().String(), "CUSTOMER"))
	r.GET("/v1/shop/:shop_id", handler.GetShop)

	mockUsecase.On("GetShopByID", shopIDstr).Return(nil, errors.New("usecase error")).Once()

	rr := performHandlerTestRequest(t, r, http.MethodGet, "/v1/shop/"+shopIDstr, nil, nil)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockUsecase.AssertExpectations(t)
}
