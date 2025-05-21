package repository

import (
	"database/sql"
	"myproject/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func GetUserByCredentials(db *sql.DB, username, password string) (*model.User, error) {
	var user model.User
	var hashedPassword string

	row := db.QueryRow("SELECT id, username, password, role FROM users WHERE username=$1", username)
	if err := row.Scan(&user.ID, &user.Username, &hashedPassword, &user.Role); err != nil {
		return nil, err
	}

	// Сравниваем хеш пароля
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(db *sql.DB, username, password string) error {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, 'user')",
		username, string(hashedPassword))
	return err
}
