package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/session"
)

func AdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		username, _ := sess.Values["username"].(string)
		role, _ := sess.Values["role"].(string)

		log.Printf("Попытка доступа к админ-панели. Пользователь: %s, Роль: %s", username, role)

		if role != "admin" {
			log.Printf("Доступ запрещен для пользователя %s с ролью %s", username, role)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		log.Printf("Доступ разрешен для администратора %s", username)

		// Получаем все записи
		bookings, err := repository.GetAllBookings(db)
		if err != nil {
			log.Printf("Ошибка получения записей: %v", err)
			http.Error(w, "Ошибка получения записей", http.StatusInternalServerError)
			return
		}

		// Получаем список пользователей
		users, err := repository.GetAllUsers(db)
		if err != nil {
			log.Printf("Ошибка получения пользователей: %v", err)
			http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
			return
		}

		// Получаем список услуг
		services, err := repository.GetServices(db)
		if err != nil {
			log.Printf("Ошибка получения услуг: %v", err)
			http.Error(w, "Ошибка получения услуг", http.StatusInternalServerError)
			return
		}

		data := &model.PageData{
			Title:    "Админ-панель - Салон красоты",
			Username: username,
			Role:     role,
			Bookings: bookings,
			Users:    users,
			Services: services,
		}
		RenderTemplate(w, "admin.html", data, true)
	}
}

func DeleteBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод запроса: %s", r.Method)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		bookingIDStr := r.FormValue("booking_id")
		if bookingIDStr == "" {
			log.Printf("Отсутствует ID записи в запросе")
			http.Error(w, "Отсутствует ID записи", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(bookingIDStr)
		if err != nil {
			log.Printf("Ошибка преобразования ID записи: %v", err)
			http.Error(w, "Неверный формат ID записи", http.StatusBadRequest)
			return
		}

		log.Printf("Попытка удаления записи с ID: %d", id)
		err = repository.DeleteBooking(db, id)
		if err != nil {
			log.Printf("Ошибка при удалении записи: %v", err)
			http.Error(w, "Ошибка при удалении записи", http.StatusInternalServerError)
			return
		}

		log.Printf("Запись с ID %d успешно удалена", id)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func DeleteUserBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		id, err := strconv.Atoi(r.FormValue("booking_id"))
		if err != nil {
			http.Error(w, "Неверный ID записи", http.StatusBadRequest)
			return
		}

		err = repository.DeleteBooking(db, id)
		if err != nil {
			http.Error(w, "Ошибка при удалении записи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

func EditBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод запроса: %s", r.Method)
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		role, _ := sess.Values["role"].(string)
		if role != "admin" {
			log.Printf("Попытка редактирования записи не администратором")
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		bookingIDStr := r.FormValue("booking_id")
		newBookingTimeStr := r.FormValue("new_booking_time")

		if bookingIDStr == "" || newBookingTimeStr == "" {
			log.Printf("Отсутствуют необходимые параметры")
			http.Error(w, "Отсутствуют необходимые параметры", http.StatusBadRequest)
			return
		}

		bookingID, err := strconv.Atoi(bookingIDStr)
		if err != nil {
			log.Printf("Ошибка преобразования ID записи: %v", err)
			http.Error(w, "Неверный формат ID записи", http.StatusBadRequest)
			return
		}

		// Парсим новое время записи
		newBookingTime, err := time.Parse("2006-01-02T15:04", newBookingTimeStr)
		if err != nil {
			log.Printf("Ошибка парсинга времени записи: %v", err)
			http.Error(w, "Неверный формат времени", http.StatusBadRequest)
			return
		}

		log.Printf("Попытка обновления записи с ID: %d на время: %v", bookingID, newBookingTime)

		err = repository.UpdateBookingTime(db, bookingID, newBookingTime)
		if err != nil {
			log.Printf("Ошибка при обновлении времени записи: %v", err)
			http.Error(w, "Ошибка при обновлении времени записи", http.StatusInternalServerError)
			return
		}

		log.Printf("Время записи успешно обновлено для записи с ID: %d", bookingID)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

// CreateBookingHandler обрабатывает создание новой записи администратором
func CreateBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Проверяем права администратора
		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		role, _ := sess.Values["role"].(string)
		if role != "admin" {
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		// Получаем данные из формы
		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "Не указано имя пользователя", http.StatusBadRequest)
			return
		}

		serviceID, err := strconv.Atoi(r.FormValue("service_id"))
		if err != nil {
			http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
			return
		}

		bookingTime, err := time.Parse("2006-01-02T15:04", r.FormValue("booking_time"))
		if err != nil {
			http.Error(w, "Неверный формат времени", http.StatusBadRequest)
			return
		}

		// Создаем запись
		booking := model.Booking{
			Username:    username,
			ServiceID:   serviceID,
			BookingTime: bookingTime.Format("15:04, 02.01.2006"),
		}
		err = repository.CreateBooking(db, booking)
		if err != nil {
			log.Printf("Ошибка создания записи: %v", err)
			http.Error(w, "Ошибка создания записи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

// UpdateBookingStatusHandler обрабатывает изменение статуса записи
func UpdateBookingStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Неподдерживаемый метод запроса: %s", r.Method)
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		role, _ := sess.Values["role"].(string)
		if role != "admin" {
			log.Printf("Попытка изменения статуса записи не администратором")
			http.Error(w, "Доступ запрещен", http.StatusForbidden)
			return
		}

		bookingIDStr := r.FormValue("booking_id")
		newStatus := r.FormValue("status")

		if bookingIDStr == "" || newStatus == "" {
			log.Printf("Отсутствуют необходимые параметры")
			http.Error(w, "Отсутствуют необходимые параметры", http.StatusBadRequest)
			return
		}

		bookingID, err := strconv.Atoi(bookingIDStr)
		if err != nil {
			log.Printf("Ошибка преобразования ID записи: %v", err)
			http.Error(w, "Неверный формат ID записи", http.StatusBadRequest)
			return
		}

		// Проверяем валидность статуса
		validStatuses := map[string]bool{
			"pending":   true,
			"confirmed": true,
			"cancelled": true,
			"completed": true,
		}
		if !validStatuses[newStatus] {
			log.Printf("Неверный статус записи: %s", newStatus)
			http.Error(w, "Неверный статус записи", http.StatusBadRequest)
			return
		}

		log.Printf("Попытка обновления статуса записи с ID: %d на статус: %s", bookingID, newStatus)

		err = repository.UpdateBookingStatus(db, bookingID, newStatus)
		if err != nil {
			log.Printf("Ошибка при обновлении статуса записи: %v", err)
			http.Error(w, "Ошибка при обновлении статуса записи", http.StatusInternalServerError)
			return
		}

		log.Printf("Статус записи успешно обновлен для записи с ID: %d", bookingID)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
