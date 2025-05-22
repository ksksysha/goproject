package handler

import (
	"database/sql"
	"log"
	"mygoproject/internal/model"
	"mygoproject/internal/repository"
	"mygoproject/internal/session"
	"net/http"
)

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

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
