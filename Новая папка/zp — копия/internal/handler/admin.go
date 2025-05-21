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

		bookings, err := repository.GetAllBookings(db)
		if err != nil {
			log.Printf("Ошибка получения записей: %v", err)
			http.Error(w, "Ошибка получения записей", http.StatusInternalServerError)
			return
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

func DeleteBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		id, _ := strconv.Atoi(r.FormValue("id"))
		_ = repository.DeleteBooking(db, id)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func DeleteUserBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		id, _ := strconv.Atoi(r.FormValue("id"))
		_ = repository.DeleteBooking(db, id)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
