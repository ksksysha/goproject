package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	log.Println("Инициализация подключения к базе данных...")
	connStr := "host=localhost port=5432 user=postgres password=0000 dbname=myproject sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	log.Println("Проверка соединения с базой данных...")
	if err = db.Ping(); err != nil {
		log.Fatalf("База данных не отвечает: %v", err)
	}

	log.Println("✅ База данных успешно подключена!")
	return db
}
