package domain

import "time"

type Repository interface {
	GetUserPassword(username string) (string, error)
	CreateUser(username, password string) error
	GetUserRole(username string) (string, error)
	GetUserID(username string) (int, error)
	GetServices() ([]Service, error)
	GetUserBookings(userID int) ([]Booking, error)
	GetAllBookings() ([]AdminBooking, error)
	CreateBooking(userID, serviceID int, bookingTime time.Time) error
	DeleteBooking(bookingID int) error
	DeleteUserBooking(bookingID, userID int) error
}
