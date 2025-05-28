package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Store provides a wrapper around the sqlc generated Queries
// and the underlying database connection.
type Store struct {
	*Queries // Embeds the sqlc-generated Queries struct
	db       *sql.DB
}

// NewStore creates a new Store instance.
// It takes an *sql.DB connection and initializes the sqlc Queries.
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db), // 'New' here refers to sqlc's generated New function in db.go
		db:      db,
	}
}

// InitDB connects to the PostgreSQL database using environment variables
// and returns a *Store instance.
func InitDB() (*Store, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Printf("Attempting to connect to database: %s:%s/%s", dbHost, dbPort, dbName)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(10)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	// Ping the database with retry logic
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		err = dbConn.Ping()
		if err == nil {
			log.Println("Successfully connected to the database!")
			return NewStore(dbConn), nil // Return the custom Store
		}
		log.Printf("Database not ready (attempt %d/%d): %v. Retrying in 2 seconds...", i+1, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
}

// CloseDB closes the underlying database connection.
func (s *Store) CloseDB() {
	if s.db != nil {
		log.Println("Closing database connection.")
		err := s.db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}
