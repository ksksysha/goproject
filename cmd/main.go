package main

import (
	"log"
	"net/http"

	"myproject/internal/config"
	"myproject/internal/delivery/http/handler"
	"myproject/internal/repository/postgres"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func main() {
	// Инициализация конфигурации
	cfg := config.LoadConfig()

	// Инициализация репозитория
	_, err := postgres.NewPostgresRepository(cfg.DBConnStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Инициализация обработчиков
	h, err := handler.NewHandler(store)
	if err != nil {
		log.Fatal("Ошибка инициализации обработчиков:", err)
	}

	// Настройка маршрутов
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", h.HomePageHandler)
	http.HandleFunc("/home", h.HomeHandler)
	http.HandleFunc("/about", h.AboutHandler)
	http.HandleFunc("/contacts", h.ContactsHandler)
	http.HandleFunc("/services", h.ServicesHandler)
	http.HandleFunc("/login", h.LoginHandler)
	http.HandleFunc("/register", h.RegisterHandler)
	http.HandleFunc("/profile", h.ProfileHandler)
	http.HandleFunc("/logout", h.LogoutHandler)
	http.HandleFunc("/book", h.BookServiceHandler)
	http.HandleFunc("/admin", h.AdminHandler)
	http.HandleFunc("/delete-booking", h.DeleteBookingHandler)
	http.HandleFunc("/delete-user-booking", h.DeleteUserBookingHandler)
	http.HandleFunc("/404", h.NotFoundHandler)

	log.Println("Запуск сервера на http://localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
