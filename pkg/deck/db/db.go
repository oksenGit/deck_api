package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Init initializes the database connection.
func Init() *sql.DB {
	var err error

	// Build connection string
	connStr := os.Getenv("DB_URL")

	// Connect to PostgreSQL database
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the database connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	return DB
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
