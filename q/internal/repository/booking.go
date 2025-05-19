package repository

import (
	"database/sql"
	"myproject/internal/model"
)

func GetAllBookings(db *sql.DB) ([]model.Booking, error) {
	rows, err := db.Query("SELECT id, service_id, date FROM bookings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(&b.ID, &b.ServiceID, &b.Date); err != nil {
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
	_, err := db.Exec("INSERT INTO bookings (username, service_id, date) VALUES ($1, $2, $3)",
		booking.Username, booking.ServiceID, booking.Date)
	return err
}

func GetUserBookings(db *sql.DB, username string) ([]model.Booking, error) {
	rows, err := db.Query("SELECT id, service_id, date FROM bookings WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(&b.ID, &b.ServiceID, &b.Date); err != nil {
			return nil, err
		}
		b.Username = username
		bookings = append(bookings, b)
	}
	return bookings, nil
}
