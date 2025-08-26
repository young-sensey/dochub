package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/database"
	"backend/routes"
)

// setupIntegrationTestDB создает тестовую БД для интеграционных тестов
func setupIntegrationTestDB(t *testing.T) *sql.DB {
	db, err := database.Connect()
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to database: %v", err)
		return nil
	}

	// Очищаем таблицу для тестов
	_, err = db.Exec("DELETE FROM documents")
	if err != nil {
		t.Fatalf("Failed to clean test table: %v", err)
	}

	return db
}

// TestFullCRUDWorkflow тестирует полный цикл CRUD операций
func TestFullCRUDWorkflow(t *testing.T) {
	db := setupIntegrationTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	// Создаем роутер
	router := routes.SetupRoutes(db)

	// 1. Создание документа
	createData := map[string]string{
		"title":   "Integration Test Doc",
		"content": "Integration test content",
		"author":  "Integration Test Author",
	}
	jsonData, _ := json.Marshal(createData)

	req := httptest.NewRequest("POST", "/dock", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var createdDoc map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createdDoc)
	docID := int(createdDoc["id"].(float64))

	// Получить все документы и вывести их
	req = httptest.NewRequest("GET", "/dock", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	t.Logf("All docs: %s", w.Body.String())

	// 2. Получение списка документов
	req = httptest.NewRequest("GET", "/dock", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var documents []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &documents)

	if len(documents) == 0 {
		t.Error("Expected at least one document after creation")
	}

	// 3. Получение конкретного документа
	req = httptest.NewRequest("GET", fmt.Sprintf("/dock/%d", docID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var retrievedDoc map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &retrievedDoc)

	if retrievedDoc["title"] != createData["title"] {
		t.Errorf("Expected title %s, got %s", createData["title"], retrievedDoc["title"])
	}

	// 4. Обновление документа
	updateData := map[string]string{
		"title":   "Updated Integration Test Doc",
		"content": "Updated integration test content",
		"author":  "Updated Integration Test Author",
	}
	jsonData, _ = json.Marshal(updateData)

	req = httptest.NewRequest("PUT", fmt.Sprintf("/dock/%d", docID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var updatedDoc map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updatedDoc)

	if updatedDoc["title"] != updateData["title"] {
		t.Errorf("Expected updated title %s, got %s", updateData["title"], updatedDoc["title"])
	}

	// 5. Удаление документа
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/dock/%d", docID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// 6. Проверка удаления
	req = httptest.NewRequest("GET", fmt.Sprintf("/dock/%d", docID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 after deletion, got %d", w.Code)
	}
}

// TestHealthEndpoint тестирует health check endpoint
func TestHealthEndpoint(t *testing.T) {
	db := setupIntegrationTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	router := routes.SetupRoutes(db)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var healthResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &healthResponse)

	if healthResponse["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", healthResponse["status"])
	}
}

// TestInvalidRequests тестирует обработку некорректных запросов
func TestInvalidRequests(t *testing.T) {
	db := setupIntegrationTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	router := routes.SetupRoutes(db)

	// Тест некорректного JSON
	req := httptest.NewRequest("POST", "/dock", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}

	// Тест несуществующего документа
	req = httptest.NewRequest("GET", "/dock/999999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent document, got %d", w.Code)
	}

	// Тест некорректного ID
	req = httptest.NewRequest("GET", "/dock/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
	}
}
