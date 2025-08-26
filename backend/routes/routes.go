package routes

import (
	"database/sql"
	"net/http"

	"backend/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes настраивает все маршруты API
func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Создаем обработчики
	docHandler := handlers.NewDocumentHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)

	// Маршруты для документов
	r.HandleFunc("/dock", docHandler.GetDocuments).Methods("GET")
	r.HandleFunc("/dock", docHandler.CreateDocument).Methods("POST")
	r.HandleFunc("/dock/{id}", docHandler.GetDocument).Methods("GET")
	r.HandleFunc("/dock/{id}", docHandler.UpdateDocument).Methods("PUT")
	r.HandleFunc("/dock/{id}", docHandler.DeleteDocument).Methods("DELETE")
	r.HandleFunc("/dock/{id}/download", docHandler.DownloadDocument).Methods("GET")

	// Маршруты для категорий
	r.HandleFunc("/categories", categoryHandler.GetCategories).Methods("GET")
	r.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	r.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	r.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	return r
}
