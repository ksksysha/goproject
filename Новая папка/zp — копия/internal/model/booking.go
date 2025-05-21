// internal/model/booking.go
package model

type Service struct {
	ID    int
	Name  string
	Price float64
}

type Booking struct {
	ID        int
	Username  string
	ServiceID int
	Date      string
	Service   Service
}
