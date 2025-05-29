package domain

import (
	"github.com/google/uuid"
	"time"
)

//type Role string

//const (
//	RoleCustomer Role = "CUSTOMER"
//	RoleSeller   Role = "SELLER"
//	RoleAdmin    Role = "ADMIN"
//)

type Customer struct {
	ID              uuid.UUID `json:"id"`
	Phone           string    `json:"phone"`
	Name            string    `json:"name"`
	Role            string    `json:"role"`
	TotalSpent      float64   `json:"total_spent"`
	CashbackBalance float64   `json:"cashback_balance"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}
