package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin"

	// Генерируем новый хеш
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Ошибка при генерации хеша: %v", err)
	}

	fmt.Printf("Новый хеш: %s\n", string(hashedPassword))
	fmt.Printf("Длина хеша: %d\n", len(hashedPassword))

	// Проверяем существующий хеш
	existingHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	err = bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(password))
	if err != nil {
		fmt.Printf("Ошибка проверки пароля: %v\n", err)
	} else {
		fmt.Println("Пароль верный!")
	}
}
