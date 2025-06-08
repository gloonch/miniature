package application

import "miniature/customer/internal/domain"

type CustomerUsecase interface {
	RegisterCustomer(phone, name, role string) (*domain.Customer, error)
	GetCustomerByID(id string) (*domain.Customer, error)
	GetCustomerByPhone(phone string) (*domain.Customer, error)
	UpdateCustomer(*domain.Customer) error
}
