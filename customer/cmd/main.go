package main

import (
	"github.com/gloonch/miniature/customer/internal/application"
	"github.com/gloonch/miniature/customer/internal/domain"
	"github.com/gloonch/miniature/customer/internal/infra/postgres"
	"github.com/gloonch/miniature/customer/internal/interfaces"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	db := postgres.NewPostgresConnection()
	defer db.Close()

	// Dependency Injection
	var repo domain.CustomerRepository = postgres.NewCustomerRepository(db)
	var usecase application.CustomerUsecase = application.NewCustomerService(repo)
	handler := interfaces.NewCustomerHandler(usecase)

	// Routing
	http.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.RegisterCustomer(w, r)
			return
		}
		if r.Method == http.MethodGet {
			handler.GetCustomerByPhone(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	addr := "localhost:8080"
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Customer service running on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
