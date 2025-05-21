package repository

import (
	"database/sql"
	"myproject/internal/model"
)

func GetAllBookings(db *sql.DB) ([]model.Booking, error) {
	query := `
		SELECT b.id, b.service_id, b.booking_time, u.username
		FROM bookings b
		JOIN users u ON b.user_id = u.id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(&b.ID, &b.ServiceID, &b.Date, &b.Username); err != nil {
			return nil, err
		}
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

	_, err = db.Exec("INSERT INTO bookings (user_id, service_id, booking_time) VALUES ($1, $2, $3)",
		userID, booking.ServiceID, booking.Date)
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
		if err := rows.Scan(&b.ID, &b.ServiceID, &b.Date, &serviceName, &servicePrice); err != nil {
			return nil, err
		}
		b.Username = username
		b.Service = model.Service{
			ID:    b.ServiceID,
			Name:  serviceName,
			Price: servicePrice,
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}
