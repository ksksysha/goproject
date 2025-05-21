package router

import (
	"database/sql"
	"net/http"

	"myproject/internal/handler"
	"myproject/internal/middleware"
)

func SetupRoutes(db *sql.DB) {
	// Статические файлы
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("zp/assets"))))

	// API эндпоинты
	http.Handle("/api/services/", handler.GetServicesByCategoryHandler(db))

	// Публичные маршруты
	http.HandleFunc("/", handler.HomePageHandler(db))
	http.HandleFunc("/login", handler.LoginHandler(db))
	http.HandleFunc("/register", handler.RegisterHandler(db))
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/services", handler.ServicesHandler(db))
	http.HandleFunc("/services/nails", handler.ServicesCategoryHandler(db, "nails"))
	http.HandleFunc("/services/hair", handler.ServicesCategoryHandler(db, "hair"))
	http.HandleFunc("/services/lashes", handler.ServicesCategoryHandler(db, "lashes"))
	http.HandleFunc("/services/brows", handler.ServicesCategoryHandler(db, "brows"))
	http.HandleFunc("/services/description/", handler.ServiceDescriptionHandler(db))
	http.HandleFunc("/about", handler.AboutHandler)
	http.HandleFunc("/contacts", handler.ContactsHandler)

	// Защищенные маршруты
	http.HandleFunc("/profile", middleware.AuthMiddleware(handler.ProfileHandler(db)))
	http.HandleFunc("/book", middleware.AuthMiddleware(handler.BookServiceHandler(db)))
	http.HandleFunc("/delete-user-booking", middleware.AuthMiddleware(handler.DeleteUserBookingHandler(db)))

	// Административные маршруты
	http.HandleFunc("/admin", middleware.AdminMiddleware(handler.AdminHandler(db)))
	http.HandleFunc("/admin/bookings", middleware.AdminMiddleware(handler.AdminBookingsHandler(db)))
	http.HandleFunc("/admin/delete-booking", middleware.AdminMiddleware(handler.DeleteBookingHandler(db)))
}
