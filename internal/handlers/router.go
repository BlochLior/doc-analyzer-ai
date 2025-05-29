package handlers

import (
	"net/http"

	"github.com/BlochLior/doc-analyzer-ai/internal/db"
)

type Handlers struct {
	Store *db.Store
}

func NewHandlers(store *db.Store) *Handlers {
	return &Handlers{Store: store}
}

// SetupRouter configures and returns the main HTTP router.
func (h *Handlers) SetupRouter() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", h.healthCheckHandler) // Call handler method on h

	// Document-related routes
	mux.HandleFunc("POST /documents", h.createDocumentHandler) // Specific method for POST
	mux.HandleFunc("GET /documents", h.listDocumentsHandler)   // Specific method for GET
	// You might add specific handlers for GET /documents/{id}, PUT /documents/{id}, DELETE /documents/{id} later

	return mux
}

// healthCheckHandler is a simple handler to check if the server is alive.
func (h *Handlers) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
