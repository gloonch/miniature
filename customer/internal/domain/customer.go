package domain

import "time"

//type Role string

//const (
//	RoleCustomer Role = "CUSTOMER"
//	RoleSeller   Role = "SELLER"
//	RoleAdmin    Role = "ADMIN"
//)

type Customer struct {
	ID              string
	Phone           string
	Name            string
	Role            string
	TotalSpent      float64
	CashbackBalance float64
	CreatedAt       time.Time
}
