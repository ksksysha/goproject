package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
)

func AdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := session.Store.Get(r, "session")
		role, _ := sess.Values["role"].(string)
		if role != "admin" {
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		bookings, err := repository.GetAllBookings(db)
		if err != nil {
			http.Error(w, "Ошибка получения записей", http.StatusInternalServerError)
			return
		}

		data := &model.PageData{Bookings: bookings}
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
