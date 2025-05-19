package http

import (
	"html/template"
	"net/http"

	"myproject/internal/domain"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := domain.PageData{Title: "Вход"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		valid, err := h.authUC.Login(username, password)
		if err != nil {
			data.ErrorMessage = "Неверный логин или пароль"
		} else if valid {
			session, _ := h.store.Get(r, "session-name")
			session.Values["username"] = username
			session.Save(r, w)

			admin, err := h.authUC.IsAdmin(username)
			if err != nil {
				data.ErrorMessage = "Ошибка при проверке роли"
			} else if admin {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			} else {
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			}
		} else {
			data.ErrorMessage = "Неверный логин или пароль"
		}
	}

	data.Content = template.HTML(`...`) // HTML как в исходном коде
	h.renderTemplate(w, "login.html", &data, true)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := domain.PageData{Title: "Регистрация"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := h.authUC.Register(username, password)
		if err != nil {
			data.ErrorMessage = "Ошибка при регистрации: пользователь уже существует"
		} else {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
	}

	data.Content = template.HTML(`...`) // HTML как в исходном коде
	h.renderTemplate(w, "register.html", &data, true)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	delete(session.Values, "username")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
