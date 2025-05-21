// internal/model/booking.go
package model

import (
	"time"
)

type Service struct {
	ID    int
	Name  string
	Price float64
}

type Booking struct {
	ID          int
	Username    string
	ServiceID   int
	BookingTime string
	Service     Service
	IsExpired   bool
	Status      string // "pending", "confirmed", "cancelled", "completed"
}

// CalculateExpired вычисляет и устанавливает поле IsExpired
func (b *Booking) CalculateExpired() {
	if b.BookingTime == "" {
		b.IsExpired = false
		return
	}

	// Пробуем оба формата времени
	bookingTime, err := time.Parse("2006-01-02 15:04:05", b.BookingTime)
	if err != nil {
		bookingTime, err = time.Parse("15:04, 02.01.2006", b.BookingTime)
		if err != nil {
			b.IsExpired = false
			return
		}
	}

	b.IsExpired = bookingTime.Before(time.Now())
}
