package config

import (
	"database/sql"
	"fmt"
	"os"

	"myproject/internal/repository/migrations"

	_ "github.com/lib/pq"
)

// InitDB инициализирует подключение к базе данных и применяет миграции
func InitDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения к базе данных: %v", err)
	}

	// Применяем миграции
	if err := migrations.ApplyMigrations(db); err != nil {
		return nil, fmt.Errorf("ошибка применения миграций: %v", err)
	}

	return db, nil
}
