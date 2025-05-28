package main

import (
	"log"
	"net/http"
	"os" // Don't forget to import os for env var

	"github.com/BlochLior/doc-analyzer-ai/internal/db"
	"github.com/BlochLior/doc-analyzer-ai/internal/handlers"
)

func main() {
	// Initialize database connection and get the custom Store
	store, err := db.InitDB() // Calls the InitDB function from internal/db/store.go
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.CloseDB() // Ensure the database connection is closed when main exits

	// IMPORTANT: You'll eventually pass this 'store' object to your handlers
	// so they can interact with the database.
	// For example:
	// h := handlers.New(store)
	// router := h.SetupRouter() // If handlers used the store

	// For now, continue with the basic router setup
	router := handlers.SetupRouter()

	log.Printf("Go server starting on port %s...", os.Getenv("PORT"))
	// Listen on port 8080 (as exposed by Docker)
	log.Fatal(http.ListenAndServe(":8080", router))
}
