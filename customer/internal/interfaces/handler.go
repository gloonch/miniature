package interfaces

import (
	"github.com/gloonch/miniature/customer/internal/application"
	//"github.com/gloonch/miniature/customer/internal/domain"
	"encoding/json"
	"net/http"
)

type CustomerHandler struct {
	usecase application.CustomerUsecase
}

func NewCustomerHandler(u application.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: u}
}

func (h *CustomerHandler) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	customer, err := h.usecase.RegisterCustomer(req.Phone, req.Name, req.Role)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) GetCustomerByPhone(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		http.Error(w, "Phone required", http.StatusBadRequest)
		return
	}
	customer, err := h.usecase.GetCustomerByPhone(phone)
	if err != nil {
		http.Error(w, "Error fetching customer", http.StatusInternalServerError)
		return
	}
	if customer == nil {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, customer)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
