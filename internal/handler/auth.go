package handler

import (
	"database/sql"
	"log"
	"net/http"

	"mygoproject/internal/model"
	"mygoproject/internal/repository"
	"mygoproject/internal/session"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			data := &model.PageData{
				Title: "Вход - Салон красоты",
			}
			RenderTemplate(w, "login.html", data, true)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		log.Printf("Попытка входа пользователя: %s", username)

		user, err := repository.GetUserByCredentials(db, username, password)
		if err != nil || user == nil {
			log.Printf("Ошибка аутентификации для пользователя %s: %v", username, err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		log.Printf("Успешная аутентификация пользователя %s с ролью %s", username, user.Role)

		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии при входе: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		sess.Values["username"] = user.Username
		sess.Values["role"] = user.Role
		sess.Values[session.UserIDKey] = user.ID
		err = sess.Save(r, w)
		if err != nil {
			log.Printf("Ошибка сохранения сессии: %v", err)
			http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
			return
		}

		log.Printf("Сессия сохранена, перенаправление пользователя %s на %s", username, user.Role)

		if user.Role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	}
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			data := &model.PageData{
				Title: "Регистрация - Салон красоты",
			}
			RenderTemplate(w, "register.html", data, true)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		err := repository.CreateUser(db, username, password)
		if err != nil {
			http.Error(w, "Ошибка регистрации", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.Store.Get(r, "session-name")
	delete(sess.Values, "username")
	delete(sess.Values, "role")
	delete(sess.Values, session.UserIDKey)
	sess.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
