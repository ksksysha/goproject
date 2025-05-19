package session

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("your-secret-key"))

func Init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
}
