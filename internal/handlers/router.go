package handlers

import (
	"log"
	"net/http"
)

// SetupRouter configures and returns the main HTTP router with middleware
func SetupRouter() http.Handler {
	// Create a new ServeMux (Go's built-in HTTP request multiplexer)
	mux := http.NewServeMux()

	// --- Define application routes here ---
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/analyze", analyzeDocumentHandler)
	mux.HandleFunc("/documents", listDocumentsHandler)

	log.Println("HTTP router configured (without CORS).")
	return mux
}

// healthCheckHandler is a simple handler to check if the server is alive
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// placeholder handlers for now:
func analyzeDocumentHandler(w http.ResponseWriter, r *http.Request) {
	// This will later contain logic to process document uploads/text
	// and call the AI service.
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Analyze Document endpoint not implemented yet."))
}

func listDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	// This will later contain logic to list stored documents.
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("List Documents endpoint not implemented yet."))
}
