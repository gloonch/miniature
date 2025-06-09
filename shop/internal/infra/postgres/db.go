package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

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
