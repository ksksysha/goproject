package main

import (
	"log"

	"mygoproject/internal/config"
	"mygoproject/internal/repository"
)

func main() {
	// Инициализация подключения к базе данных
	db := config.InitDB()
	defer db.Close()

	// Создание нового пользователя-администратора
	err := repository.CreateAdmin(db, "admin", "admin")
	if err != nil {
		log.Fatalf("Ошибка создания администратора: %v", err)
	}

	log.Println("Администратор успешно создан")
}
