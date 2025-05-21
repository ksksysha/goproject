package handler

import (
	"database/sql"
	"net/http"

	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
)

// HomePageHandler обрабатывает запросы к главной странице
func HomePageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		sess, _ := session.Store.Get(r, "session-name")
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

		renderTemplate(w, "home.html", data, true)
	}
}

// AboutHandler обрабатывает запросы к странице "О нас"
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := model.PageData{
		Title: "О нас",
	}

	renderTemplate(w, "about.html", &data, true)
}

// ContactsHandler обрабатывает запросы к странице контактов
func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	data := model.PageData{
		Title: "Контакты",
	}

	renderTemplate(w, "contacts.html", &data, true)
}
