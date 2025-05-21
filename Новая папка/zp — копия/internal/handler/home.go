package handler

import (
	"database/sql"
	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
	"net/http"
)

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := session.Store.Get(r, "session")
		username, _ := sess.Values["username"].(string)
		role, _ := sess.Values["role"].(string)

		services, err := repository.GetServices(db)
		if err != nil {
			http.Error(w, "Ошибка загрузки услуг", http.StatusInternalServerError)
			return
		}

		data := &model.PageData{
			Title:    "Главная - Салон красоты",
			Username: username,
			Role:     role,
			Services: services,
		}
		RenderTemplate(w, "layout.html", data, true)
	}
}
