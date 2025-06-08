package interfaces

type CreateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Role  string `json:"role" binding:"required"`
}

type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
}
