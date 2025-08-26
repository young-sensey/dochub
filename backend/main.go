package main

import (
	"log"
	"net/http"

	"backend/database"
	"backend/routes"
)

func main() {
	// Подключение к базе данных
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Создание таблиц
	err = database.CreateTables(db)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Настройка маршрутов
	r := routes.SetupRoutes(db)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
