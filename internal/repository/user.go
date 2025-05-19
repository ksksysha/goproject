package repository

import (
	"database/sql"
	"myproject/internal/model"
)

func GetUserByCredentials(db *sql.DB, username, password string) (*model.User, error) {
	row := db.QueryRow("SELECT id, username, role FROM users WHERE username=$1 AND password=$2", username, password)

	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(db *sql.DB, username, password string) error {
	_, err := db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, 'user')", username, password)
	return err
}
