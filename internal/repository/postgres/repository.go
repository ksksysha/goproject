package postgres

import (
	"database/sql"
	"time"

	"myproject/internal/domain"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) GetUserPassword(username string) (string, error) {
	var password string
	err := r.db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&password)
	return password, err
}

func (r *PostgresRepository) CreateUser(username, password string) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	return err
}

func (r *PostgresRepository) GetUserRole(username string) (string, error) {
	var role string
	err := r.db.QueryRow("SELECT role FROM users WHERE username=$1", username).Scan(&role)
	return role, err
}

func (r *PostgresRepository) GetUserID(username string) (int, error) {
	var userID int
	err := r.db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userID)
	return userID, err
}

func (r *PostgresRepository) GetServices() ([]domain.Service, error) {
	rows, err := r.db.Query("SELECT id, name, price FROM services")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []domain.Service
	for rows.Next() {
		var service domain.Service
		if err := rows.Scan(&service.ID, &service.Name, &service.Price); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *PostgresRepository) GetUserBookings(userID int) ([]domain.Booking, error) {
	rows, err := r.db.Query(`
		SELECT b.id, s.name, b.booking_time
		FROM bookings b
		JOIN services s ON b.service_id = s.id
		WHERE b.user_id = $1
		ORDER BY b.booking_time ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		var bookingTime time.Time
		if err := rows.Scan(&booking.ID, &booking.ServiceName, &bookingTime); err != nil {
			return nil, err
		}
		booking.BookingTime = bookingTime.Format("2006-01-02 15:04")
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (r *PostgresRepository) GetAllBookings() ([]domain.AdminBooking, error) {
	rows, err := r.db.Query(`
		SELECT b.id, u.username, s.name, b.booking_time
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		JOIN services s ON b.service_id = s.id
		ORDER BY b.booking_time ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.AdminBooking
	for rows.Next() {
		var booking domain.AdminBooking
		var bookingTime time.Time
		if err := rows.Scan(&booking.ID, &booking.Username, &booking.ServiceName, &bookingTime); err != nil {
			return nil, err
		}
		booking.IsExpired = bookingTime.Before(time.Now())
		booking.BookingTime = bookingTime.Format("2006-01-02 15:04")
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (r *PostgresRepository) CreateBooking(userID, serviceID int, bookingTime time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO bookings (user_id, service_id, booking_time) VALUES ($1, $2, $3)",
		userID, serviceID, bookingTime,
	)
	return err
}

func (r *PostgresRepository) DeleteBooking(bookingID int) error {
	_, err := r.db.Exec("DELETE FROM bookings WHERE id = $1", bookingID)
	return err
}

func (r *PostgresRepository) DeleteUserBooking(bookingID, userID int) error {
	_, err := r.db.Exec("DELETE FROM bookings WHERE id = $1 AND user_id = $2", bookingID, userID)
	return err
}
