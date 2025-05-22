// internal/model/booking.go
package model

import (
	"log"
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

	bookingTime, err := time.Parse("15:04, 02.01.2006", b.BookingTime)
	if err != nil {
		log.Printf("Ошибка парсинга времени в CalculateExpired: %v", err)
		b.IsExpired = false
		return
	}

	b.IsExpired = bookingTime.Before(time.Now())
}
