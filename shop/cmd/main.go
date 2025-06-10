package main

import (
	"log"
	// "os"

	// "github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"

	// "github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/application"
	// "github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/infra/postgres"
	// "github.com/segment-sources/sources-backend-takehome-assignment/shop/internal/interfaces"
	// "github.com/segment-sources/sources-backend-takehome-assignment/pkg/db" // Assuming a shared db package
)

func main() {
	// Load environment variables
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found, using environment variables")
	// }

	// Database connection (Commented out for now)
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// dbSSLMode := os.Getenv("DB_SSLMODE")

	// database, err := db.NewPostgresDB(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer database.Close()

	// log.Println("Database connected successfully")

	// Initialize repository (Commented out for now)
	// shopRepo := postgres.NewShopRepository(database.GetDB()) // Assuming NewShopRepository takes *sql.DB

	// Initialize service (Commented out for now)
	// shopService := application.NewShopService(shopRepo)

	// Initialize handler (Commented out for now)
	// shopHandler := interfaces.NewShopHandler(shopService)

	// Initialize router (Commented out for now)
	// router := interfaces.NewRouter(shopHandler)

	// Start server (Commented out for now)
	// port := os.Getenv("SHOP_SERVICE_PORT")
	// if port == "" {
	// 	port = "8081" // Default port for shop service
	// }
	// log.Printf("Shop service starting on port %s", port)
	// if err := router.Run(":" + port); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }

	log.Println("Shop service main function executed (placeholders active).")
}
