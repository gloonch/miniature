package application

import (
	"database/sql"
	"errors"
	"miniature/product/internal/domain"
	"time"

	"github.com/google/uuid"
)

type productService struct {
	repo                 domain.Repository
	shopOwnershipChecker domain.ShopOwnershipCheckerRepository // Added
}

func NewProductService(repo domain.Repository, shopChecker domain.ShopOwnershipCheckerRepository) Usecase { // Updated
	return &productService{repo: repo, shopOwnershipChecker: shopChecker} // Updated
}

func (s *productService) CreateProduct(shopIDStr, name, description string, price float64, sku string, stockQuantity int, creatingUserIDStr string) (*domain.Product, error) {
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return nil, errors.New("invalid shop_id format")
	}

	// Authorization: Verify creatingUserIDStr owns shopID.
	isOwner, err := s.shopOwnershipChecker.IsShopOwner(creatingUserIDStr, shopIDStr)
	if err != nil {
		// Log the error for server-side insight
		// log.Printf("Error checking shop ownership for user %s, shop %s: %v", creatingUserIDStr, shopIDStr, err)
		return nil, errors.New("could not verify shop ownership") // Generic error to client
	}
	if !isOwner {
		return nil, errors.New("user not authorized to add products to this shop")
	}

	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if stockQuantity < 0 {
		return nil, errors.New("stock quantity cannot be negative")
	}

	product := &domain.Product{
		ID:            uuid.New(),
		ShopID:        shopID,
		Name:          name,
		Description:   description,
		Price:         price,
		SKU:           sku,
		StockQuantity: stockQuantity,
		IsActive:      true, // Default to active
		CreatedAt:     time.Now(),
	}

	err = s.repo.Create(product)
	if err != nil {
		// Consider specific error for SKU unique violation (e.g. check pq.Error as before)
		return nil, err
	}
	return product, nil
}

func (s *productService) GetProductByID(id string) (*domain.Product, error) {
	// TODO - Authorization: Consider if any user can fetch any product by ID, or if there are restrictions.
	// For now, open access if product exists.
	return s.repo.FindByID(id)
}

func (s *productService) GetProductsByShopID(shopIDStr string /*, requestingUserIDStr string */) ([]*domain.Product, error) {
	// _, err := uuid.Parse(shopIDStr) // Validate shopIDStr if needed, though repo will handle bad UUIDs from DB side.
	// if err != nil {
	// 	return nil, errors.New("invalid shop_id format")
	// }

	// TODO - Authorization: Consider if any user can fetch products for any shop.
	// Or if it should be restricted (e.g., only shop owner, or if shop is public).
	// For now, open access.
	return s.repo.FindByShopID(shopIDStr)
}

func (s *productService) UpdateProduct(
	productIDStr string,
	name *string,
	description *string,
	price *float64,
	sku *string,
	stockQuantity *int,
	isActive *bool,
	requestingUserIDStr string,
) (*domain.Product, error) {
	product, err := s.repo.FindByID(productIDStr)
	if err != nil {
		return nil, errors.New("database error while finding product: " + err.Error())
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Authorization: Verify requestingUserIDStr owns product.ShopID.
	isOwner, err := s.shopOwnershipChecker.IsShopOwner(requestingUserIDStr, product.ShopID.String())
	if err != nil {
		return nil, errors.New("could not verify shop ownership for product update")
	}
	if !isOwner {
		return nil, errors.New("user not authorized to update this product")
	}

	if name != nil {
		product.Name = *name
	}
	if description != nil {
		product.Description = *description
	}
	if price != nil {
		if *price < 0 {
			return nil, errors.New("price cannot be negative")
		}
		product.Price = *price
	}
	if sku != nil {
		product.SKU = *sku
	}
	if stockQuantity != nil {
		if *stockQuantity < 0 {
			return nil, errors.New("stock quantity cannot be negative")
		}
		product.StockQuantity = *stockQuantity
	}
	if isActive != nil {
		product.IsActive = *isActive
	}

	err = s.repo.Update(product)
	if err != nil {
		// Handle specific errors like SKU conflict if necessary
		// if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		// 	 if strings.Contains(pqErr.Constraint, "uq_shop_sku") {
		//		 return nil, errors.New("product with this SKU already exists in this shop")
		//	 }
		//}
		return nil, errors.New("database error while updating product: " + err.Error())
	}
	return product, nil
}

func (s *productService) DeleteProduct(productIDStr string, requestingUserIDStr string) error {
	product, err := s.repo.FindByID(productIDStr)
	if err != nil && err != sql.ErrNoRows { // If it's a real DB error, not just not found
		return errors.New("database error while finding product: " + err.Error())
	}
	if product == nil || err == sql.ErrNoRows { // Product doesn't exist
		return sql.ErrNoRows // Propagate not found
	}

	// Authorization: Verify requestingUserIDStr owns product.ShopID.
	isOwner, err := s.shopOwnershipChecker.IsShopOwner(requestingUserIDStr, product.ShopID.String())
	if err != nil {
		return errors.New("could not verify shop ownership for product deletion")
	}
	if !isOwner {
		return errors.New("user not authorized to delete this product")
	}
	return s.repo.Delete(productIDStr)
}

// Add placeholders for other service methods
// ...
