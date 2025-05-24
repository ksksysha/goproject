package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	log.Println("Инициализация подключения к базе данных...")

	config := GetDBConfig()
	user, password := GetDBCredentials()
	if user == "" || password == "" {
		log.Fatal("DB_USER или DB_PASSWORD не установлены в secrets.env или переменных окружения")
	}

	db, err := sql.Open("postgres", config.GetConnectionString(user, password))
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
