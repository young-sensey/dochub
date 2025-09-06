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

	// Создаем таблицу users
	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(usersQuery)
	if err != nil {
		return err
	}

	// Миграции: заменить author на user_id, заполнить существующие записи user_id=2
	// 1) Добавить столбец user_id, если его нет
	_, _ = db.Exec(`ALTER TABLE documents ADD COLUMN IF NOT EXISTS user_id INTEGER`)
	// 2) Проставить user_id=2 для пустых значений
	_, _ = db.Exec(`UPDATE documents SET user_id = 2 WHERE user_id IS NULL`)
	// 3) Удалить столбец author, если существует
	_, _ = db.Exec(`ALTER TABLE documents DROP COLUMN IF EXISTS author`)
	// 4) Добавить ограничение NOT NULL
	_, _ = db.Exec(`ALTER TABLE documents ALTER COLUMN user_id SET NOT NULL`)
	// 5) Добавить внешний ключ на users, если отсутствует
	_, _ = db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints
				WHERE constraint_type = 'FOREIGN KEY'
				  AND table_name = 'documents'
				  AND constraint_name = 'documents_user_id_fkey'
			) THEN
				ALTER TABLE documents
				ADD CONSTRAINT documents_user_id_fkey FOREIGN KEY (user_id)
				REFERENCES users(id) ON DELETE CASCADE;
			END IF;
		END$$;`)

	log.Println("Tables created successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
