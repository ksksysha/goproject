package domain

import (
	"html/template"
)

type PageData struct {
	Title        string
	Content      template.HTML
	ErrorMessage string
	Username     string
	Services     []Service
	Bookings     []Booking
}

type Service struct {
	ID    int
	Name  string
	Price float64
}

type Booking struct {
	ID          int
	ServiceName string
	BookingTime string
}

type AdminBooking struct {
	ID          int
	Username    string
	ServiceName string
	BookingTime string
	IsExpired   bool
}
