package main

import (
	_ "github.com/lib/pq"
	"log"
	"miniature/product/internal/application"
	"miniature/product/internal/infra/postgres"
	"miniature/product/internal/interfaces"
)

func main() {

	db := postgres.NewPostgresConnection()
	defer db.Close()

	repo := postgres.NewRepository(db)
	shopRepo := postgres.NewShopRepository(db)
	usecase := application.NewProductService(repo, shopRepo)
	productHandler := interfaces.NewHandler(usecase)
	route := interfaces.NewRouter(productHandler)

	addr := "localhost:8082"
	route.Run(addr)

	log.Printf("Product service running on: %s", addr)
}
