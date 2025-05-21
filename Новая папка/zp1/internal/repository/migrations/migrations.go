package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// ApplyMigrations применяет все миграции из директории migrations
func ApplyMigrations(db *sql.DB) error {
	// Получаем список файлов миграций
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return fmt.Errorf("ошибка чтения директории миграций: %v", err)
	}

	// Сортируем файлы по имени (они должны быть в формате XXX_name.sql)
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Создаем таблицу для отслеживания примененных миграций, если её нет
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы migrations: %v", err)
	}

	// Применяем каждую миграцию
	for _, file := range migrationFiles {
		// Проверяем, была ли миграция уже применена
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = $1)", file).Scan(&exists)
		if err != nil {
			return fmt.Errorf("ошибка проверки миграции %s: %v", file, err)
		}
		if exists {
			continue
		}

		// Читаем содержимое файла миграции
		content, err := ioutil.ReadFile(filepath.Join(".", file))
		if err != nil {
			return fmt.Errorf("ошибка чтения файла миграции %s: %v", file, err)
		}

		// Начинаем транзакцию
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("ошибка начала транзакции для %s: %v", file, err)
		}

		// Выполняем миграцию
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка выполнения миграции %s: %v", file, err)
		}

		// Отмечаем миграцию как примененную
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", file)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка отметки миграции %s как примененной: %v", file, err)
		}

		// Завершаем транзакцию
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("ошибка завершения транзакции для %s: %v", file, err)
		}

		fmt.Printf("Применена миграция: %s\n", file)
	}

	return nil
}
