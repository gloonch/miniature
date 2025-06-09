package main

import (
	"log"
	"miniature/shop/internal/application"
	"miniature/shop/internal/infra/postgres"
	"miniature/shop/internal/interfaces"
)

func main() {

	db := postgres.NewPostgresConnection()
	defer db.Close()

	repo := postgres.NewPostgresShopRepository(db)
	service := application.NewShopService(repo)
	handler := interfaces.NewShopHandler(service)
	route := interfaces.NewRouter(handler)

	addr := "localhost:8081"
	route.Run(addr)

	log.Printf("Shop service running on: %s", addr)
}
