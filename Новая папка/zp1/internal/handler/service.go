package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"myproject/internal/repository"
)

func GetServicesByCategoryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Получаем slug категории из URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 4 {
			http.Error(w, "Неверный формат URL", http.StatusBadRequest)
			return
		}
		categorySlug := parts[3]

		// Получаем услуги для категории
		services, err := repository.GetServicesByCategory(db, categorySlug)
		if err != nil {
			http.Error(w, "Ошибка при получении услуг", http.StatusInternalServerError)
			return
		}

		// Отправляем ответ в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	}
}
