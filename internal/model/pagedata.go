package model

type PageData struct {
	Title             string
	Username          string
	Role              string
	UserID            int
	Services          []Service
	Users             []User
	Bookings          []Booking
	Content           string
	SelectedServiceID string
}
