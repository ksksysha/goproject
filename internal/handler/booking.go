package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
)

func BookServiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		sess, _ := session.Store.Get(r, "session")
		username, _ := sess.Values["username"].(string)

		serviceID, err := strconv.Atoi(r.FormValue("service"))
		if err != nil {
			http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
			return
		}
		date := r.FormValue("date")

		booking := model.Booking{
			Username:  username,
			ServiceID: serviceID,
			Date:      date,
		}

		err = repository.CreateBooking(db, booking)
		if err != nil {
			http.Error(w, "Ошибка записи", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
