package postgres

import (
	"database/sql"
	"github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/domain"
)

type postgresShopRepository struct {
	db *sql.DB
}

func NewPostgresShopRepository(db *sql.DB) domain.ShopRepository {
	return &postgresShopRepository{db: db}
}

func (r *postgresShopRepository) Create(shop *domain.Shop) error {
	query := `INSERT INTO shops (id, name, owner_id, address, is_active, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, shop.ID, shop.Name, shop.OwnerID, shop.Address, shop.IsActive, shop.CreatedAt, shop.UpdatedAt)
	return err
}

import (
	"database/sql"
	"github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/domain"
)

type postgresShopRepository struct {
	db *sql.DB
}

func NewPostgresShopRepository(db *sql.DB) domain.ShopRepository {
	return &postgresShopRepository{db: db}
}

func (r *postgresShopRepository) Create(shop *domain.Shop) error {
	query := `INSERT INTO shops (id, name, owner_id, address, is_active, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, shop.ID, shop.Name, shop.OwnerID, shop.Address, shop.IsActive, shop.CreatedAt, shop.UpdatedAt)
	return err
}

func (r *postgresShopRepository) FindByID(id string) (*domain.Shop, error) {
	shop := &domain.Shop{}
	query := `SELECT id, name, owner_id, address, is_active, created_at, updated_at
              FROM shops WHERE id = $1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&shop.ID, &shop.Name, &shop.OwnerID, &shop.Address, &shop.IsActive, &shop.CreatedAt, &shop.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a domain-specific error like domain.ErrShopNotFound
		}
		return nil, err
	}
	return shop, nil
}

func (r *postgresShopRepository) FindByOwnerID(ownerID string) ([]*domain.Shop, error) {
	var shops []*domain.Shop
	query := `SELECT id, name, owner_id, address, is_active, created_at, updated_at
              FROM shops WHERE owner_id = $1`
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		shop := &domain.Shop{}
		err := rows.Scan(&shop.ID, &shop.Name, &shop.OwnerID, &shop.Address, &shop.IsActive, &shop.CreatedAt, &shop.UpdatedAt)
		if err != nil {
			return nil, err // Or collect errors and continue
		}
		shops = append(shops, shop)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *postgresShopRepository) Update(shop *domain.Shop) error {
	query := `UPDATE shops SET name = $1, address = $2, is_active = $3, updated_at = $4
              WHERE id = $5`
	_, err := r.db.Exec(query, shop.Name, shop.Address, shop.IsActive, shop.UpdatedAt, shop.ID)
	return err
}

func (r *postgresShopRepository) Delete(id string) error {
	query := `DELETE FROM shops WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Or a custom domain error for not found
	}
	return nil
}
