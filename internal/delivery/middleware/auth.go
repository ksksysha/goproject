package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AuthMiddleware(store *sessions.CookieStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func AdminMiddleware(store *sessions.CookieStore, authUC *usecase.AuthUseCase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		admin, err := authUC.IsAdmin(username)
		if err != nil || !admin {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
