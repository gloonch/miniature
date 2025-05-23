package interfaces

type CreateCustomerRequest struct {
	Phone string `json:"phone"`
	Name  string `json:"name"`
	Role  string `json:"role"` // Optional
}
