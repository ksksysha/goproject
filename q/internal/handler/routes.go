package handler

import (
	"database/sql"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Статические файлы
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Обработчики страниц
	mux.HandleFunc("/", PageHandler)
	mux.HandleFunc("/home", PageHandler)
	mux.HandleFunc("/about", PageHandler)
	mux.HandleFunc("/services", PageHandler)
	mux.HandleFunc("/contacts", PageHandler)

	// Обработчики действий
	mux.HandleFunc("/login", LoginHandler(db))
	mux.HandleFunc("/register", RegisterHandler(db))
	mux.HandleFunc("/profile", ProfileHandler(db))
	mux.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/book", BookServiceHandler(db))
	mux.HandleFunc("/admin", AdminHandler(db))
	mux.HandleFunc("/delete-booking", DeleteBookingHandler(db))
	mux.HandleFunc("/delete-user-booking", DeleteUserBookingHandler(db))

	// Обработчик 404
	mux.HandleFunc("/404", NotFoundHandler)
}
