package usecase

import (
	"myproject/internal/domain"
	"time"
)

type BookingUseCase struct {
	repo domain.Repository
}

func NewBookingUseCase(repo domain.Repository) *BookingUseCase {
	return &BookingUseCase{repo: repo}
}

func (uc *BookingUseCase) GetServices() ([]domain.Service, error) {
	return uc.repo.GetServices()
}

func (uc *BookingUseCase) GetUserBookings(username string) ([]domain.Booking, error) {
	userID, err := uc.repo.GetUserID(username)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetUserBookings(userID)
}

func (uc *BookingUseCase) BookService(username string, serviceID int, bookingTime time.Time) error {
	userID, err := uc.repo.GetUserID(username)
	if err != nil {
		return err
	}
	return uc.repo.CreateBooking(userID, serviceID, bookingTime)
}

func (uc *BookingUseCase) DeleteUserBooking(username string, bookingID int) error {
	userID, err := uc.repo.GetUserID(username)
	if err != nil {
		return err
	}
	return uc.repo.DeleteUserBooking(bookingID, userID)
}
