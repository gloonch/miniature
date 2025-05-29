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

	customerTable := "customers"
	customerDDL := `
		CREATE TABLE customers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone VARCHAR(20) UNIQUE NOT NULL,
			name TEXT,
			role TEXT DEFAULT 'OWNER',
			total_spent NUMERIC DEFAULT 0,
			cashback_balance NUMERIC DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
	customerSeed := `
		INSERT INTO customers (id, phone, name, role)
		VALUES 
			(gen_random_uuid(), '09121234567', 'Alice', 'OWNER'),
			(gen_random_uuid(), '09121234568', 'Bob', 'CUSTOMER'),
			(gen_random_uuid(), '09121234569', 'Charlie', 'SELLER');`

	_ = ensureTableExists(db, customerTable, customerDDL)
	_ = seedTable(db, customerTable, customerSeed)

	shopTable := "shops"

	shopDDL := `
	CREATE TABLE shops (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		owner_id UUID REFERENCES customers (id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		description TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	shopSeed := `
	INSERT INTO shops (id, owner_id, name, description)
	VALUES 
		(gen_random_uuid(), 
		 (SELECT id FROM customers WHERE phone = '09121234567' LIMIT 1), 
		 'فروشگاه آلیس', 'پوشاک زنانه'),

		(gen_random_uuid(), 
		 (SELECT id FROM customers WHERE phone = '09121234568' LIMIT 1), 
		 'فروشگاه باب', 'اکسسوری و زیورآلات');`

	_ = ensureTableExists(db, shopTable, shopDDL)
	_ = seedTable(db, shopTable, shopSeed)

	productTable := "products"

	productDDL := `
	CREATE TABLE products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		code VARCHAR(50) UNIQUE NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		price NUMERIC NOT NULL,
		image_url TEXT,
		stock INTEGER DEFAULT 0,
		category TEXT,
		shop_id UUID REFERENCES shops (id) ON DELETE CASCADE,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	productSeed := `
	INSERT INTO products (id, code, title, description, price, stock, category, shop_id)
	VALUES 
		(gen_random_uuid(), 'SL38', 'شلوار جین زنانه', 'شلوار جین آبی سایز 38', 490000, 10, 'پوشاک',
		 (SELECT id FROM shops WHERE name = 'فروشگاه آلیس' LIMIT 1)),

		(gen_random_uuid(), 'BLK01', 'کیف چرمی مشکی', 'کیف رودوشی چرمی کلاسیک', 750000, 5, 'اکسسوری',
		 (SELECT id FROM shops WHERE name = 'فروشگاه باب' LIMIT 1));`

	_ = ensureTableExists(db, productTable, productDDL)
	_ = seedTable(db, productTable, productSeed)

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

func ensureTableExists(db *sql.DB, tableName string, createStmt string) error {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = $1
		);`
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking if table %s exists: %w", tableName, err)
	}

	if exists {
		log.Printf("✅ Table '%s' already exists.\n", tableName)
		return nil
	}

	log.Printf("⚠️ Table '%s' does not exist. Creating...\n", tableName)
	_, err = db.Exec(createStmt)
	if err != nil {
		return fmt.Errorf("creating table %s: %w", tableName, err)
	}

	log.Printf("✅ Table '%s' created successfully.\n", tableName)
	return nil
}

func seedTable(db *sql.DB, tableName string, insertStmt string) error {
	var count int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	if err != nil {
		return fmt.Errorf("checking table %s: %w", tableName, err)
	}

	if count > 0 {
		log.Printf("ℹ️ Table '%s' already has data.\n", tableName)
		return nil
	}

	log.Printf("⏳ Seeding data into '%s'...\n", tableName)
	_, err = db.Exec(insertStmt)
	if err != nil {
		return fmt.Errorf("inserting data into %s: %w", tableName, err)
	}

	log.Printf("✅ Seeded data into '%s'.\n", tableName)
	return nil
}
