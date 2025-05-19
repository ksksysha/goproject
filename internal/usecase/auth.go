package usecase

import (
	"myproject/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo domain.Repository
}

func NewAuthUseCase(repo domain.Repository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

func (uc *AuthUseCase) Login(username, password string) (bool, error) {
	dbPassword, err := uc.repo.GetUserPassword(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	return err == nil, err
}

func (uc *AuthUseCase) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return uc.repo.CreateUser(username, string(hashedPassword))
}

func (uc *AuthUseCase) IsAdmin(username string) (bool, error) {
	role, err := uc.repo.GetUserRole(username)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}
