package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"backend/entities"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	db *sql.DB
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetCategories возвращает список всех категорий
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []entities.Category
	for rows.Next() {
		var category entities.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// CreateCategory создает новую категорию
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req entities.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO categories (name, description) 
	VALUES ($1, $2) 
	RETURNING id, name, description, created_at, updated_at`

	var category entities.Category
	err := h.db.QueryRow(query, req.Name, req.Description).
		Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GetCategory возвращает категорию по ID
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var category entities.Category
	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1"
	err = h.db.QueryRow(query, id).
		Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory обновляет категорию по ID
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req entities.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
	UPDATE categories 
	SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP 
	WHERE id = $3 
	RETURNING id, name, description, created_at, updated_at`

	var category entities.Category
	err = h.db.QueryRow(query, req.Name, req.Description, id).
		Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory удаляет категорию по ID
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
