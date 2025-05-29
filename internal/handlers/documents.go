package handlers

import (
	"database/sql" // Required for sql.NullString
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BlochLior/doc-analyzer-ai/internal/db"
	"github.com/google/uuid" // For generating UUIDs if not passed in
)

// CreateDocumentRequest represents the request body for creating a new document.
type CreateDocumentRequest struct {
	Title   string `json:"title"` // New: title field
	Content string `json:"content"`
	Summary string `json:"summary"`
	AIModel string `json:"ai_model"` // ai_model can be an empty string for nullable
}

// DocumentResponse represents the structure for a single document in API responses.
// We'll use this for both creation responses and list responses.
type DocumentResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"` // New: title field
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	AIModel   *string   `json:"ai_model,omitempty"` // Pointer for nullable, omitempty for JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListDocumentsResponse represents the response body for listing documents.
type ListDocumentsResponse struct {
	Documents []DocumentResponse `json:"documents"`
}

// createDocumentHandler handles the creation of a new document.
// POST /documents
func (h *Handlers) createDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Prepare parameters for sqlc's CreateDocument method
	params := db.CreateDocumentParams{
		Title:   req.Title,
		Content: req.Content,
		Summary: req.Summary,
	}

	// Handle the nullable ai_model field. sqlc generates sql.NullString for TEXT NULLABLE.
	if req.AIModel != "" {
		params.AiModel = sql.NullString{String: req.AIModel, Valid: true}
	} else {
		params.AiModel = sql.NullString{Valid: false} // Explicitly set as invalid if empty
	}

	// Pass the request context
	document, err := h.Store.CreateDocument(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		log.Printf("Error creating document in DB: %v", err)
		return
	}

	// Map the sqlc generated document to our public response struct
	resp := DocumentResponse{
		ID:        document.ID,
		Title:     document.Title,
		Content:   document.Content,
		Summary:   document.Summary,
		CreatedAt: document.CreatedAt,
		UpdatedAt: document.UpdatedAt,
	}
	// Handle nullable AIModel for response
	if document.AiModel.Valid {
		resp.AIModel = &document.AiModel.String
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(resp)
	log.Printf("Document created with ID: %s", document.ID)
}

// listDocumentsHandler handles listing all documents.
// GET /documents
func (h *Handlers) listDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	documents, err := h.Store.ListDocuments(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve documents", http.StatusInternalServerError)
		log.Printf("Error listing documents from DB: %v", err)
		return
	}

	var respDocuments []DocumentResponse
	for _, doc := range documents {
		docResp := DocumentResponse{
			ID:        doc.ID,
			Title:     doc.Title,
			Content:   doc.Content,
			Summary:   doc.Summary,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		}
		if doc.AiModel.Valid {
			docResp.AIModel = &doc.AiModel.String
		}
		respDocuments = append(respDocuments, docResp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListDocumentsResponse{Documents: respDocuments})
	log.Println("Listed all documents.")
}
