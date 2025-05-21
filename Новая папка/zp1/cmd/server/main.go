package main

import (
	"log"
	"net/http"
	"os"

	"myproject/internal/config"
	"myproject/internal/router"
	"myproject/internal/session"
)

func main() {
	// Устанавливаем переменные окружения по умолчанию, если они не заданы
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "0000")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "myproject")
	}

	// Инициализируем базу данных
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	// Инициализируем сессии
	session.Init()

	// Настраиваем маршруты
	router.SetupRoutes(db)

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
