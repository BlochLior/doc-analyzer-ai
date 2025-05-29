package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BlochLior/doc-analyzer-ai/internal/db"
	"github.com/BlochLior/doc-analyzer-ai/internal/handlers"
)

func main() {
	// Initialize database connection and get the custom Store
	store, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.CloseDB() // Ensure the database connection is closed when main exits

	// Create a new Handlers instance and pass the database store to it
	apiHandlers := handlers.NewHandlers(store)

	// Setup the main HTTP router using the Handlers instance
	// CORRECTED LINE BELOW: Call SetupRouter as a method on apiHandlers
	router := apiHandlers.SetupRouter()

	log.Printf("Go server starting on port %s...", os.Getenv("PORT")) // os.Getenv("PORT") will be empty string, but :8080 will work
	log.Fatal(http.ListenAndServe(":8000", router))
}
