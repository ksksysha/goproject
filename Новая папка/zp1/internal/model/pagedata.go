package model

type PageData struct {
	Title             string
	Username          string
	Role              string
	UserID            int
	Services          []Service
	Categories        []Category
	CurrentCategory   string
	Bookings          []Booking
	Content           string
	SelectedServiceID string
}
