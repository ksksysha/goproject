package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"mygoproject/internal/model"
	"mygoproject/internal/repository"
	"mygoproject/internal/session"
)

func BookServiceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		sess, err := session.Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			http.Error(w, "Ошибка сессии", http.StatusInternalServerError)
			return
		}

		username, ok := sess.Values["username"].(string)
		if !ok {
			log.Printf("Ошибка получения username из сессии")
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}

		serviceID, err := strconv.Atoi(r.FormValue("service_id"))
		if err != nil {
			log.Printf("Ошибка преобразования service_id: %v", err)
			http.Error(w, "Неверный ID услуги", http.StatusBadRequest)
			return
		}

		bookingTime := r.FormValue("booking_time")
		if bookingTime == "" {
			log.Printf("Пустое время записи")
			http.Error(w, "Необходимо указать время записи", http.StatusBadRequest)
			return
		}

		// Преобразуем время из формата datetime-local в нужный формат
		parsedTime, err := time.Parse("2006-01-02T15:04", bookingTime)
		if err != nil {
			log.Printf("Ошибка парсинга времени: %v", err)
			http.Error(w, "Неверный формат времени записи", http.StatusBadRequest)
			return
		}

		log.Printf("Попытка создания записи: username=%s, service_id=%d, booking_time=%s",
			username, serviceID, parsedTime.Format("15:04, 02.01.2006"))

		booking := model.Booking{
			Username:    username,
			ServiceID:   serviceID,
			BookingTime: parsedTime.Format("15:04, 02.01.2006"),
		}

		err = repository.CreateBooking(db, booking)
		if err != nil {
			log.Printf("Ошибка создания записи: %v", err)
			http.Error(w, "Ошибка записи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Запись успешно создана для пользователя %s", username)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
