package application

import (
	"github.com/gloonch/miniature/customer/internal/domain"
	"github.com/google/uuid"
	"time"
)

type customerService struct {
	repo domain.CustomerRepository
}

func NewCustomerService(repo domain.CustomerRepository) *customerService {
	return &customerService{repo: repo}
}

func (cs *customerService) RegisterCustomer(phone, name, role string) (*domain.Customer, error) {
	customer := domain.Customer{
		ID:              uuid.NewString(),
		Phone:           phone,
		Name:            name,
		Role:            role,
		TotalSpent:      0,
		CashbackBalance: 0,
		CreatedAt:       time.Time{},
	}

	err := cs.repo.Create(&customer)

	return &customer, err
}

func (cs *customerService) GetCustomerByID(id string) (*domain.Customer, error) {
	return cs.repo.FindByID(id)
}

func (cs *customerService) GetCustomerByPhone(phone string) (*domain.Customer, error) {
	return cs.repo.FindByPhone(phone)
}

func (cs *customerService) UpdateCustomer(customer *domain.Customer) error {
	return cs.repo.Update(customer)
}
