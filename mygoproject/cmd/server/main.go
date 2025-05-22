package main

import (
	"log"
	"net/http"
	"os"

	"mygoproject/internal/config"
	"mygoproject/internal/handler"
	"mygoproject/internal/session"
)

func main() {
	// Настраиваем логгер
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Запуск сервера...")

	// Инициализация базы данных
	log.Println("Подключение к базе данных...")
	db := config.InitDB()
	defer db.Close()

	session.Init() // Инициализация сессий

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, db)

	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
