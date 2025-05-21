package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
)

// AdminHandler обрабатывает запросы к административной панели
func AdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		username, _ := sess.Values["username"].(string)
		role, _ := sess.Values["role"].(string)

		log.Printf("Попытка доступа к админ-панели. Пользователь: %s, Роль: %s", username, role)

		if role != "admin" {
			log.Printf("Доступ запрещен для пользователя %s с ролью %s", username, role)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		log.Printf("Доступ разрешен для администратора %s", username)

		// Получаем все записи
		bookings, err := repository.GetAllBookings(db)
		if err != nil {
			log.Printf("Ошибка получения записей: %v", err)
			http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
			return
		}

		// Вычисляем статус истекших записей
		for i := range bookings {
			bookings[i].CalculateExpired()
		}

		data := &model.PageData{
			Title:    "Админ-панель - Салон красоты",
			Username: username,
			Role:     role,
			Bookings: bookings,
		}
		RenderTemplate(w, "admin.html", data, true)
	}
}

// AdminBookingsHandler обрабатывает запросы к странице управления записями
func AdminBookingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем все записи
		bookings, err := repository.GetAllBookings(db)
		if err != nil {
			http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
			return
		}

		// Вычисляем статус истекших записей
		for i := range bookings {
			bookings[i].CalculateExpired()
		}

		data := model.PageData{
			Title:    "Управление записями",
			Bookings: bookings,
		}

		renderTemplate(w, "admin_bookings.html", &data, true)
	}
}

// DeleteBookingHandler обрабатывает запросы на удаление записи
func DeleteBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		bookingIDStr := r.FormValue("booking_id")
		if bookingIDStr == "" {
			http.Error(w, "Не указан ID записи", http.StatusBadRequest)
			return
		}

		bookingID, err := strconv.Atoi(bookingIDStr)
		if err != nil {
			http.Error(w, "Неверный формат ID записи", http.StatusBadRequest)
			return
		}

		err = repository.DeleteBooking(db, bookingID)
		if err != nil {
			http.Error(w, "Ошибка при удалении записи", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
	}
}

func DeleteUserBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		id, err := strconv.Atoi(r.FormValue("booking_id"))
		if err != nil {
			http.Error(w, "Неверный ID записи", http.StatusBadRequest)
			return
		}

		err = repository.DeleteBooking(db, id)
		if err != nil {
			http.Error(w, "Ошибка при удалении записи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
