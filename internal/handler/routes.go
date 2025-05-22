package handler

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"runtime"
)

var assetsDir string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	assetsDir = filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(currentFile))), "assets")
}

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Статические файлы
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))

	// Обработчики действий
	mux.HandleFunc("/login", LoginHandler(db))
	mux.HandleFunc("/register", RegisterHandler(db))
	mux.HandleFunc("/profile", ProfileHandler(db))
	mux.HandleFunc("/logout", LogoutHandler)
	mux.HandleFunc("/book", BookServiceHandler(db))
	mux.HandleFunc("/admin", AdminHandler(db))
	mux.HandleFunc("/delete-booking", DeleteBookingHandler(db))
	mux.HandleFunc("/delete-user-booking", DeleteUserBookingHandler(db))
	mux.HandleFunc("/admin/edit-booking", EditBookingHandler(db))
	mux.HandleFunc("/admin/create-booking", CreateBookingHandler(db))
	mux.HandleFunc("/admin/update-status", UpdateBookingStatusHandler(db))

	// Обработчик 404
	mux.HandleFunc("/404", NotFoundHandler)

	// Обработчики страниц (кроме корневого пути)
	mux.HandleFunc("/home", PageHandler)
	mux.HandleFunc("/about", PageHandler)
	mux.HandleFunc("/services", PageHandler)
	mux.HandleFunc("/contacts", PageHandler)

	// Обработчик по умолчанию для всех необработанных запросов
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/home" {
			NotFoundHandler(w, r)
			return
		}
		PageHandler(w, r)
	})
}
