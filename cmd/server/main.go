package main

import (
	"log"
	"net/http"

	"myproject/internal/config"
	"myproject/internal/handler"
	"myproject/internal/session"
)

func main() {
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
