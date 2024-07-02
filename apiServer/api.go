package apiServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/linexjlin/TRewind/chromaManager"
)

type Document struct {
	DocumentID      string `json:"document_id,omitempty"`
	DocumentName    string `json:"document_name,omitempty"`
	DocumentContent string `json:"document_content,omitempty"`
	Extra           string `json:"extra,omitempty"`
}

type ApiServer struct {
	db *chromaManager.ChromaManager
}

func NewServer(db *chromaManager.ChromaManager) *ApiServer {
	return &ApiServer{db: db}
}

func (s *ApiServer) ListenAndServe(addr string) error {
	http.HandleFunc("/", s.corsMiddleware(s.router))
	log.Printf("Server starting on %s...\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *ApiServer) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (s *ApiServer) router(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	collection := parts[1]
	action := parts[2]

	switch {
	case r.Method == "POST" && action == "upload_document":
		s.uploadDocument(w, r, collection)
	case r.Method == "DELETE" && action == "delete_document_by_id":
		s.deleteDocumentByID(w, r, collection)
	case r.Method == "DELETE" && action == "delete_document_by_name":
		s.deleteDocumentByName(w, r, collection)
	case r.Method == "GET" && action == "retrieve_document":
		s.retrieveDocument(w, r, collection)
	case r.Method == "GET" && action == "search":
		s.search(w, r, collection)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *ApiServer) uploadDocument(w http.ResponseWriter, r *http.Request, collection string) {
	var document Document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if document.DocumentName == "" {
		http.Error(w, "Document name is required", http.StatusBadRequest)
		return
	}

	id := md5Hash(document.DocumentName)
	text := fmt.Sprintf("Document name: %s\n%s", document.DocumentName, document.DocumentContent)
	metadata := map[string]string{
		"extra":    document.Extra,
		"update":   time.Now().Format("20060102150405"),
		"filename": document.DocumentName,
		"content":  document.DocumentContent,
	}

	err := s.db.UpsertDoc(collection, text, id, metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"message":     "Document uploaded successfully",
		"document_id": id,
	})
}

func (s *ApiServer) deleteDocumentByID(w http.ResponseWriter, r *http.Request, collection string) {
	var document Document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if document.DocumentID != "" {
		err := s.db.DeleteByID(collection, document.DocumentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Document deleted successfully",
	})
}

func (s *ApiServer) deleteDocumentByName(w http.ResponseWriter, r *http.Request, collection string) {
	// This function is not implemented in the original Python code
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Document deleted successfully",
	})
}

func (s *ApiServer) retrieveDocument(w http.ResponseWriter, r *http.Request, collection string) {
	query := r.URL.Query().Get("query")

	results, err := s.db.Search(collection, query, 3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var documents []map[string]interface{}
	for _, result := range results {
		docInfo := map[string]interface{}{
			"document_id":      result.ID,
			"document_name":    result.Metadata["filename"],
			"document_content": result.Metadata["content"],
			"update":           result.Metadata["update"],
			"extra":            result.Metadata["extra"],
			"distance":         result.Similarity,
		}
		documents = append(documents, docInfo)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"documents": documents,
	})
}

func (s *ApiServer) search(w http.ResponseWriter, r *http.Request, collection string) {
	query := r.URL.Query().Get("q")

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "http://localhost:8000"
	}

	results, err := s.db.Search(collection, query, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var documents []map[string]interface{}
	for _, result := range results {
		docInfo := map[string]interface{}{
			"document_id":      result.ID,
			"document_name":    result.Metadata["filename"],
			"document_content": result.Metadata["content"],
			"extra":            result.Metadata["extra"],
			"update":           result.Metadata["update"],
			"url":              fmt.Sprintf("%s/%s", serverAddr, result.ID),
			"distance":         result.Similarity,
		}
		documents = append(documents, docInfo)
	}

	json.NewEncoder(w).Encode(documents)
}
