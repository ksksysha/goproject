package http

import (
	"net/http"
	"strconv"
	"time"

	"myproject/internal/domain"
)

func (h *Handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	services, err := h.bookingUC.GetServices()
	if err != nil {
		http.Error(w, "Ошибка при получении услуг", http.StatusInternalServerError)
		return
	}

	bookings, err := h.bookingUC.GetUserBookings(username)
	if err != nil {
		http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
		return
	}

	data := &domain.PageData{
		Title:    "Личный кабинет",
		Username: username,
		Services: services,
		Bookings: bookings,
	}

	h.renderTemplate(w, "profile.html", data, false)
}

func (h *Handler) BookServiceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	admin, err := h.authUC.IsAdmin(username)
	if err != nil {
		http.Error(w, "Ошибка проверки роли пользователя", http.StatusInternalServerError)
		return
	}
	if admin {
		http.Error(w, "Администраторы не могут записываться на услуги", http.StatusForbidden)
		return
	}

	serviceID, err := strconv.Atoi(r.FormValue("service_id"))
	if err != nil {
		http.Error(w, "Неверный формат ID услуги", http.StatusBadRequest)
		return
	}

	bookingTime, err := time.Parse("2006-01-02T15:04", r.FormValue("booking_time"))
	if err != nil {
		http.Error(w, "Неверный формат времени записи", http.StatusBadRequest)
		return
	}

	err = h.bookingUC.BookService(username, serviceID, bookingTime)
	if err != nil {
		http.Error(w, "Ошибка при записи на услугу", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *Handler) DeleteUserBookingHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	bookingID, err := strconv.Atoi(r.FormValue("booking_id"))
	if err != nil {
		http.Error(w, "Отсутствует ID записи", http.StatusBadRequest)
		return
	}

	err = h.bookingUC.DeleteUserBooking(username, bookingID)
	if err != nil {
		http.Error(w, "Ошибка при удалении записи", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
