package handler

import (
	"database/sql"
	"net/http"

	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
)

func ProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := session.Store.Get(r, "session")
		username, _ := sess.Values["username"].(string)
		role, _ := sess.Values["role"].(string)

		bookings, err := repository.GetUserBookings(db, username)
		if err != nil {
			http.Error(w, "Ошибка загрузки записей", http.StatusInternalServerError)
			return
		}

		services, err := repository.GetServices(db)
		if err != nil {
			http.Error(w, "Ошибка загрузки услуг", http.StatusInternalServerError)
			return
		}

		data := &model.PageData{
			Title:    "Профиль - Салон красоты",
			Username: username,
			Role:     role,
			Bookings: bookings,
			Services: services,
		}
		RenderTemplate(w, "profile.html", data, true)
	}
}
