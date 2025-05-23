package domain

type CustomerRepository interface {
	Create(customer *Customer) error
	FindByID(id string) (*Customer, error)
	FindByPhone(phone string) (*Customer, error)
	Update(customer *Customer) error
}
