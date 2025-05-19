package usecase

import (
	"myproject/internal/domain"
)

type AdminUseCase struct {
	repo domain.Repository
}

func NewAdminUseCase(repo domain.Repository) *AdminUseCase {
	return &AdminUseCase{repo: repo}
}

func (uc *AdminUseCase) GetAllBookings() ([]domain.AdminBooking, error) {
	return uc.repo.GetAllBookings()
}

func (uc *AdminUseCase) DeleteBooking(bookingID int) error {
	return uc.repo.DeleteBooking(bookingID)
}
