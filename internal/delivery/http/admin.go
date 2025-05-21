package http

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"myproject/internal/domain"
)

func (h *Handler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	admin, err := h.authUC.IsAdmin(username)
	if err != nil || !admin {
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	bookings, err := h.adminUC.GetAllBookings()
	if err != nil {
		http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Bookings []domain.AdminBooking
	}{
		Title:    "Админ-панель",
		Bookings: bookings,
	}

	tmplPath := filepath.Join("templates", "admin.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func (h *Handler) DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	admin, err := h.authUC.IsAdmin(username)
	if err != nil || !admin {
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	bookingID, err := strconv.Atoi(r.FormValue("booking_id"))
	if err != nil {
		http.Error(w, "Отсутствует ID записи", http.StatusBadRequest)
		return
	}

	err = h.adminUC.DeleteBooking(bookingID)
	if err != nil {
		http.Error(w, "Ошибка при удалении записи", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
