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
		// Проверяем авторизацию
		userID := session.GetUserID(r)
		if userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Получаем данные из сессии
		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		username, ok := sess.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		role, _ := sess.Values["role"].(string)

		// Загружаем записи пользователя
		bookings, err := repository.GetUserBookings(db, username)
		if err != nil {
			http.Error(w, "Ошибка загрузки записей: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Загружаем список услуг
		services, err := repository.GetServices(db)
		if err != nil {
			http.Error(w, "Ошибка загрузки услуг: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Формируем данные для шаблона
		data := &model.PageData{
			Title:    "Профиль - Салон красоты",
			Username: username,
			Role:     role,
			UserID:   userID,
			Bookings: bookings,
			Services: services,
		}

		// Если передан service_id в URL, добавляем его в данные
		if serviceID := r.URL.Query().Get("service_id"); serviceID != "" {
			data.SelectedServiceID = serviceID
		}

		RenderTemplate(w, "profile.html", data, true)
	}
}
