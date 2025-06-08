package main

import (
	"log"
	"miniature/customer/internal/application"
	"miniature/customer/internal/domain"
	"miniature/customer/internal/infra/postgres"
	"miniature/customer/internal/interfaces"
)

func main() {

	db := postgres.NewPostgresConnection()
	defer db.Close()

	var repo domain.CustomerRepository = postgres.NewCustomerRepository(db)
	var usecase application.CustomerUsecase = application.NewCustomerService(repo)
	handler := interfaces.NewCustomerHandler(usecase)
	route := interfaces.NewRouter(*handler)

	addr := "localhost:8080"
	route.Run(addr)

	log.Printf("Customer service running on: %s", addr)
}
