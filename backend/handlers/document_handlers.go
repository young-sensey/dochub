package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"backend/entities"

	"github.com/gorilla/mux"
)

type DocumentHandler struct {
	db *sql.DB
}

func NewDocumentHandler(db *sql.DB) *DocumentHandler {
	return &DocumentHandler{db: db}
}

// GetDocuments возвращает список всех документов
func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр category_id из query string
	categoryID := r.URL.Query().Get("category_id")

	var rows *sql.Rows
	var err error

	if categoryID != "" {
		if categoryID == "null" {
			// Если category_id=null, возвращаем документы без категории
			rows, err = h.db.Query("SELECT id, title, content, author, file_path, category_id, created_at, updated_at FROM documents WHERE category_id IS NULL ORDER BY created_at DESC")
		} else {
			// Если указан category_id, фильтруем по категории
			rows, err = h.db.Query("SELECT id, title, content, author, file_path, category_id, created_at, updated_at FROM documents WHERE category_id = $1 ORDER BY created_at DESC", categoryID)
		}
	} else {
		// Иначе возвращаем все документы
		rows, err = h.db.Query("SELECT id, title, content, author, file_path, category_id, created_at, updated_at FROM documents ORDER BY created_at DESC")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var documents []entities.Document
	for rows.Next() {
		var doc entities.Document
		err := rows.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Author, &doc.FilePath, &doc.CategoryID, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		documents = append(documents, doc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

// CreateDocument с поддержкой загрузки файла
func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Проверяем Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && contentType[:19] == "multipart/form-data" {
		h.createDocumentWithFile(w, r)
		return
	}

	// Старый JSON-метод
	var req entities.CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO documents (title, content, author, category_id) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, title, content, author, file_path, category_id, created_at, updated_at`

	var doc entities.Document
	err := h.db.QueryRow(query, req.Title, req.Content, req.Author, req.CategoryID).
		Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Author, &doc.FilePath, &doc.CategoryID, &doc.CreatedAt, &doc.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

// createDocumentWithFile обрабатывает multipart/form-data
func (h *DocumentHandler) createDocumentWithFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Ошибка парсинга формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	author := r.FormValue("author")
	categoryIDStr := r.FormValue("category_id")

	var categoryID *int
	if categoryIDStr != "" {
		if id, err := strconv.Atoi(categoryIDStr); err == nil {
			categoryID = &id
		}
	}

	file, handler, err := r.FormFile("file")
	filePath := ""
	if err == nil && file != nil {
		defer file.Close()
		// Сохраняем файл в папку uploads
		uploadDir := "uploads"
		os.MkdirAll(uploadDir, 0755)
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
		fullPath := filepath.Join(uploadDir, filename)
		f, err := os.Create(fullPath)
		if err != nil {
			http.Error(w, "Ошибка сохранения файла: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		filePath = fullPath
	}

	query := `
	INSERT INTO documents (title, content, author, file_path, category_id) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, title, content, author, file_path, category_id, created_at, updated_at`

	var doc entities.Document
	err = h.db.QueryRow(query, title, content, author, filePath, categoryID).
		Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Author, &doc.FilePath, &doc.CategoryID, &doc.CreatedAt, &doc.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

// GetDocument возвращает документ по ID
func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var doc entities.Document
	query := "SELECT id, title, content, author, file_path, category_id, created_at, updated_at FROM documents WHERE id = $1"
	err = h.db.QueryRow(query, id).
		Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Author, &doc.FilePath, &doc.CategoryID, &doc.CreatedAt, &doc.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

// UpdateDocument обновляет документ по ID
func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req entities.UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
	UPDATE documents 
	SET title = $1, content = $2, author = $3, category_id = $4, updated_at = CURRENT_TIMESTAMP 
	WHERE id = $5 
	RETURNING id, title, content, author, file_path, category_id, created_at, updated_at`

	var doc entities.Document
	err = h.db.QueryRow(query, req.Title, req.Content, req.Author, req.CategoryID, id).
		Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Author, &doc.FilePath, &doc.CategoryID, &doc.CreatedAt, &doc.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

// DeleteDocument удаляет документ по ID
func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := h.db.Exec("DELETE FROM documents WHERE id = $1", id)
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
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DownloadDocument скачивает файл документа
func (h *DocumentHandler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("DownloadDocument", vars)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var filePath, title string
	err = h.db.QueryRow("SELECT file_path, title FROM documents WHERE id = $1", id).Scan(&filePath, &title)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if filePath == "" {
		http.Error(w, "Файл не загружен для этого документа", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
