package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Connect устанавливает подключение к PostgreSQL
func Connect() (*sql.DB, error) {
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "docflow")
	dbPassword := getEnv("DB_PASSWORD", "docflow_pass")
	dbName := getEnv("DB_NAME", "docflow_db")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateTables создает необходимые таблицы
func CreateTables(db *sql.DB) error {
	// Создаем таблицу categories
	categoriesQuery := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(categoriesQuery)
	if err != nil {
		return err
	}

	// Создаем таблицу documents с внешним ключом на categories
	documentsQuery := `
	CREATE TABLE IF NOT EXISTS documents (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT,
		author VARCHAR(255),
		file_path VARCHAR(255),
		category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(documentsQuery)
	if err != nil {
		return err
	}

	log.Println("Tables created successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
