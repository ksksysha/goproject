package main

import (
	"log"

	"myproject/internal/database"
	"myproject/internal/models"
)

func main() {
	// Инициализация подключения к базе данных
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	// Создание нового пользователя-администратора
	admin := &models.User{
		Username: "admin",
		Password: "admin",
		Role:     "admin",
	}

	// Хеширование пароля и сохранение в базу данных
	if err := admin.Create(db); err != nil {
		log.Fatalf("Ошибка создания администратора: %v", err)
	}

	log.Println("Администратор успешно создан")
}
