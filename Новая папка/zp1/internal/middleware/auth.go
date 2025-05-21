package middleware

import (
	"net/http"

	"myproject/internal/session"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r)
	}
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		role, ok := sess.Values["role"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
