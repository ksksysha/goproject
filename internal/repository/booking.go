package repository

import (
	"database/sql"
	"time"

	"mygoproject/internal/model"
)

func GetAllBookings(db *sql.DB) ([]model.Booking, error) {
	query := `
		SELECT b.id, b.service_id, b.booking_time, u.username, s.name as service_name, s.price, COALESCE(b.status, 'pending') as status
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		JOIN services s ON b.service_id = s.id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		var serviceName string
		var servicePrice float64
		var bookingTime time.Time
		if err := rows.Scan(&b.ID, &b.ServiceID, &bookingTime, &b.Username, &serviceName, &servicePrice, &b.Status); err != nil {
			return nil, err
		}
		b.Service = model.Service{
			ID:    b.ServiceID,
			Name:  serviceName,
			Price: servicePrice,
		}
		b.BookingTime = bookingTime.Format("15:04, 02.01.2006")
		b.CalculateExpired()
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func DeleteBooking(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM bookings WHERE id=$1", id)
	return err
}

func CreateBooking(db *sql.DB, booking model.Booking) error {
	// Сначала получаем user_id по username
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", booking.Username).Scan(&userID)
	if err != nil {
		return err
	}

	// Парсим строку времени обратно в time.Time
	bookingTime, err := time.Parse("15:04, 02.01.2006", booking.BookingTime)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO bookings (user_id, service_id, booking_time) VALUES ($1, $2, $3)",
		userID, booking.ServiceID, bookingTime)
	return err
}

func GetUserBookings(db *sql.DB, username string) ([]model.Booking, error) {
	query := `
		SELECT b.id, b.service_id, b.booking_time, s.name as service_name, s.price
		FROM bookings b
		JOIN services s ON b.service_id = s.id
		JOIN users u ON b.user_id = u.id
		WHERE u.username = $1
		ORDER BY b.booking_time DESC`

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		var serviceName string
		var servicePrice float64
		var bookingTime time.Time
		if err := rows.Scan(&b.ID, &b.ServiceID, &bookingTime, &serviceName, &servicePrice); err != nil {
			return nil, err
		}
		b.Username = username
		b.Service = model.Service{
			ID:    b.ServiceID,
			Name:  serviceName,
			Price: servicePrice,
		}
		b.BookingTime = bookingTime.Format("15:04, 02.01.2006")
		b.CalculateExpired()
		bookings = append(bookings, b)
	}
	return bookings, nil
}

// UpdateBookingTime обновляет время записи в базе данных
func UpdateBookingTime(db *sql.DB, bookingID int, newTime time.Time) error {
	query := `
		UPDATE bookings 
		SET booking_time = $1 
		WHERE id = $2
	`
	_, err := db.Exec(query, newTime, bookingID)
	return err
}

// UpdateBookingStatus обновляет статус записи в базе данных
func UpdateBookingStatus(db *sql.DB, bookingID int, status string) error {
	query := `
		UPDATE bookings 
		SET status = $1 
		WHERE id = $2
	`
	_, err := db.Exec(query, status, bookingID)
	return err
}
