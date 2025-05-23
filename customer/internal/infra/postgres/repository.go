package postgres

import (
	"database/sql"
	"errors"
	"github.com/gloonch/miniature/customer/internal/domain"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *customerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(c *domain.Customer) error {
	query := `
		INSERT INTO customer (id, name, phone, role, total_spent, cashback_balance, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		c.ID,
		c.Name,
		c.Phone,
		c.Role,
		c.TotalSpent,
		c.CashbackBalance,
		c.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *customerRepository) FindByID(id string) (*domain.Customer, error) {
	query := `SELECT * FROM customer WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var customer domain.Customer
	if err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Phone,
		&customer.Role,
		&customer.TotalSpent,
		&customer.CashbackBalance,
		&customer.CreatedAt); err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
		}
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindByPhone(phone string) (*domain.Customer, error) {
	query := `SELECT * FROM customer WHERE phone = $1`
	row := r.db.QueryRow(query, phone)

	var customer domain.Customer
	if err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Phone,
		&customer.Role,
		&customer.TotalSpent,
		&customer.CashbackBalance,
		&customer.CreatedAt); err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
		}
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(customer *domain.Customer) error {

	query := `
		UPDATE customers
		SET name = $1, phone = $2, role = $3, total_spent = $4, cashback_balance = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(query,
		customer.Name,
		customer.Phone,
		customer.Role,
		customer.TotalSpent,
		customer.CashbackBalance,
		customer.ID,
	)

	return err
}
