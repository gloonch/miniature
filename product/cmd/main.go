package main

import (
	"database/sql" // For *sql.DB
	"fmt"          // For placeholder connectToDB
	"log"
	"os" // For environment variables (example)

	"github.com/segment-sources/sources-backend-takehome-assignment/product/internal/application"
	"github.com/segment-sources/sources-backend-takehome-assignment/product/internal/domain"
	"github.com/segment-sources/sources-backend-takehome-assignment/product/internal/infra/postgres"
	"github.com/segment-sources/sources-backend-takehome-assignment/product/internal/interfaces"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Placeholder for database connection - in a real app, this would be more robust
// and likely come from a shared package or use a proper configuration.
func connectToDB() (*sql.DB, error) {
	// Example: Read from environment variables or use defaults
	// For simplicity, using hardcoded defaults (NOT FOR PRODUCTION)
	dbHost := os.Getenv("DB_HOST_PRODUCT")
	dbPort := os.Getenv("DB_PORT_PRODUCT")
	dbUser := os.Getenv("DB_USER_PRODUCT")
	dbPassword := os.Getenv("DB_PASSWORD_PRODUCT")
	dbName := os.Getenv("DB_NAME_PRODUCT")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "user"
	} // Replace with your actual user
	if dbPassword == "" {
		dbPassword = "password"
	} // Replace with your actual password
	if dbName == "" {
		dbName = "product_db"
	} // Replace with your actual DB name for product service

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close() // Close on ping error
		return nil, err
	}
	log.Println("Successfully connected to the product database!")
	return db, nil
}

func main() {
	log.Println("Product service starting...")

	// 1. Initialize DB connection
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 2. Initialize Auth Middleware
	// Assuming AuthMiddleware is defined in product/internal/interfaces/auth.go
	// and it's self-contained or configured via pkg/token which reads its own config.
	authMW := interfaces.AuthMiddleware()

	// 3. Initialize Repositories
	productRepo := postgres.NewPostgresProductRepository(db)
	// Explicitly cast to the domain interface for the service
	var domainProductRepo domain.ProductRepository = productRepo

	shopCheckerRepo := postgres.NewPostgresShopOwnershipCheckerRepository(db)
	// Explicitly cast to the domain interface for the service
	var domainShopCheckerRepo domain.ShopOwnershipCheckerRepository = shopCheckerRepo

	// 4. Initialize Service/Usecase
	productUsecase := application.NewProductService(domainProductRepo, domainShopCheckerRepo)
	var domainProductUsecase application.ProductUsecase = productUsecase

	// 5. Initialize Handler
	productHandler := interfaces.NewProductHandler(domainProductUsecase)

	// 6. Initialize Router
	// The NewRouter in product/internal/interfaces/route.go expects *ProductHandler and gin.HandlerFunc
	router := interfaces.NewRouter(productHandler, authMW)

	// 7. Run Server
	// Example port, should be configurable
	serverAddr := os.Getenv("PRODUCT_SERVICE_ADDR")
	if serverAddr == "" {
		serverAddr = "0.0.0.0:8082"
	}

	log.Printf("Product service running on: %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run product server: %v", err)
	}
}
