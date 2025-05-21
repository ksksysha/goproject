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

		// Получаем ID выбранной услуги из query параметра
		selectedServiceID := r.URL.Query().Get("service")

		// Получаем категории
		categories, err := repository.GetCategories(db)
		if err != nil {
			http.Error(w, "Ошибка при получении категорий", http.StatusInternalServerError)
			return
		}

		// Загружаем записи пользователя
		bookings, err := repository.GetUserBookings(db, username)
		if err != nil {
			http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
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
			Title:             "Личный кабинет",
			Username:          username,
			Role:              role,
			UserID:            userID,
			Categories:        categories,
			Bookings:          bookings,
			Services:          services,
			SelectedServiceID: selectedServiceID,
		}

		RenderTemplate(w, "profile.html", data, true)
	}
}
