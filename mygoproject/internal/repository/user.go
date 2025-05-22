package repository

import (
	"database/sql"
	"mygoproject/internal/model"

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

func CreateAdmin(db *sql.DB, username, password string) error {
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, 'admin')",
		username, string(hashedPassword))
	return err
}

// GetAllUsers возвращает список всех пользователей
func GetAllUsers(db *sql.DB) ([]model.User, error) {
	query := `SELECT id, username, role FROM users ORDER BY username`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
