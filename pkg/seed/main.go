package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	db := NewPostgresConnection()
	defer db.Close()

	ensureCustomersTable(db)
	seedCustomers(db)

}

func NewPostgresConnection() *sql.DB {

	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "postgres"
	dbPass := "password"
	dbName := "miniaturedb"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	return db
}

func ensureCustomersTable(db *sql.DB) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = 'customers'
		);
	`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking if table exists: %w", err)
	}

	if exists {
		log.Println("✅ Table 'customers' already exists.")
		return nil
	}

	log.Println("⚠️ Table 'customers' does not exist. Creating...")

	_, err = db.Exec(`
		CREATE TABLE customers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone VARCHAR(20) UNIQUE NOT NULL,
			name TEXT,
			role TEXT DEFAULT 'OWNER',
			total_spent NUMERIC DEFAULT 0,
			cashback_balance NUMERIC DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("creating customers table: %w", err)
	}

	log.Println("✅ Table 'customers' created.")
	return nil
}

func seedCustomers(db *sql.DB) error {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM customers`).Scan(&count)
	if err != nil {
		return fmt.Errorf("checking customers table: %w", err)
	}
	if count > 0 {
		log.Println("ℹ️ Customers already seeded.")
		return nil
	}

	_, err = db.Exec(`
		INSERT INTO customers (id, phone, name, role)
		VALUES 
			(gen_random_uuid(), '09121234567', 'Alice', 'OWNER'),
			(gen_random_uuid(), '09121234568', 'Bob', 'CUSTOMER'),
			(gen_random_uuid(), '09121234569', 'Charlie', 'SELLER');
	`)
	if err != nil {
		return fmt.Errorf("inserting mock customers: %w", err)
	}

	return nil
}
