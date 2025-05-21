// internal/model/booking.go
package model

import (
	"time"
)

// Booking представляет запись на услугу
type Booking struct {
	ID          int
	Username    string
	ServiceID   int
	BookingTime string
	Service     Service
	IsExpired   bool
}

// CalculateExpired вычисляет и устанавливает поле IsExpired
func (b *Booking) CalculateExpired() {
	if b.BookingTime == "" {
		b.IsExpired = false
		return
	}

	bookingTime, err := time.Parse("2006-01-02 15:04:05", b.BookingTime)
	if err != nil {
		b.IsExpired = false
		return
	}

	b.IsExpired = bookingTime.Before(time.Now())
}
