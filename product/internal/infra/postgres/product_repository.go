package postgres

import (
	"database/sql"
	"miniature/product/internal/domain"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(product *domain.Product) error {
	query := `INSERT INTO products
              (id, shop_id, name, description, price, sku, stock_quantity, is_active, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query,
		product.ID, product.ShopID, product.Name, product.Description, product.Price,
		product.SKU, product.StockQuantity, product.IsActive, product.CreatedAt,
	)
	return err
}

func (r *repository) FindByID(id string) (*domain.Product, error) {
	product := &domain.Product{}
	query := `SELECT id, shop_id, name, description, price, sku, stock_quantity, is_active, created_at
              FROM products WHERE id = $1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&product.ID, &product.ShopID, &product.Name, &product.Description, &product.Price,
		&product.SKU, &product.StockQuantity, &product.IsActive, &product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Standard way to indicate not found
		}
		return nil, err
	}
	return product, nil
}

func (r *repository) FindByShopID(shopID string) ([]*domain.Product, error) {
	var products []*domain.Product
	query := `SELECT id, shop_id, name, description, price, sku, stock_quantity, is_active, created_at
              FROM products WHERE shop_id = $1`
	rows, err := r.db.Query(query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID, &product.ShopID, &product.Name, &product.Description, &product.Price,
			&product.SKU, &product.StockQuantity, &product.IsActive, &product.CreatedAt,
		)
		if err != nil {
			return nil, err // Or collect errors and continue
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) Update(product *domain.Product) error {
	query := `UPDATE products SET
                name = $1,
                description = $2,
                price = $3,
                sku = $4,
                stock_quantity = $5,
                is_active = $6
              WHERE id = $7 AND shop_id = $8` // shop_id in WHERE for safety, though id is PK
	_, err := r.db.Exec(query,
		product.Name, product.Description, product.Price, product.SKU,
		product.StockQuantity, product.IsActive, product.ID, product.ShopID,
	)
	return err
}

func (r *repository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Standard way to indicate not found or nothing deleted
	}
	return nil
}
