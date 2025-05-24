package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения системы")
	}
	if err := godotenv.Overload("secrets.env"); err != nil {
		log.Println("Файл secrets.env не найден, логин и пароль берутся из переменных окружения")
	}
}

type DBConfig struct {
	Host    string
	Port    string
	DBName  string
	SSLMode string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Host:    getEnvOrDefault("DB_HOST", "localhost"),
		Port:    getEnvOrDefault("DB_PORT", "5432"),
		DBName:  getEnvOrDefault("DB_NAME", "myproject"),
		SSLMode: getEnvOrDefault("DB_SSLMODE", "disable"),
	}
}

func GetDBCredentials() (user, password string) {
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	return
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c DBConfig) GetConnectionString(user, password string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, user, password, c.DBName, c.SSLMode)
}
