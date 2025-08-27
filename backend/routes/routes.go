package routes

import (
	"database/sql"
	"net/http"

	"backend/handlers"
	"backend/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes настраивает все маршруты API
func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Создаем обработчики
	docHandler := handlers.NewDocumentHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// Публичные маршруты авторизации
	r.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Защищенные маршруты
	api := r.NewRoute().Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Маршруты для документов
	api.HandleFunc("/dock", docHandler.GetDocuments).Methods("GET")
	api.HandleFunc("/dock", docHandler.CreateDocument).Methods("POST")
	api.HandleFunc("/dock/{id}", docHandler.GetDocument).Methods("GET")
	api.HandleFunc("/dock/{id}", docHandler.UpdateDocument).Methods("PUT")
	api.HandleFunc("/dock/{id}", docHandler.DeleteDocument).Methods("DELETE")
	api.HandleFunc("/dock/{id}/download", docHandler.DownloadDocument).Methods("GET")

	// Маршруты для категорий
	api.HandleFunc("/categories", categoryHandler.GetCategories).Methods("GET")
	api.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	api.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	api.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	api.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	return r
}
