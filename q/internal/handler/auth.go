package handler

import (
	"database/sql"
	"net/http"

	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
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

		user, err := repository.GetUserByCredentials(db, username, password)
		if err != nil || user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sess, _ := session.Store.Get(r, "session")
		sess.Values["username"] = user.Username
		sess.Values["role"] = user.Role
		sess.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	sess, _ := session.Store.Get(r, "session")
	delete(sess.Values, "username")
	delete(sess.Values, "role")
	sess.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
