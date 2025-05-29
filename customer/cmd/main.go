package main

import (
	"github.com/gloonch/miniature/customer/internal/application"
	"github.com/gloonch/miniature/customer/internal/domain"
	"github.com/gloonch/miniature/customer/internal/infra/postgres"
	"github.com/gloonch/miniature/customer/internal/interfaces"
	"log"
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
