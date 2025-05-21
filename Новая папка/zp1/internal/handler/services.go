package handler

import (
	"database/sql"
	"net/http"

	"myproject/internal/model"
	"myproject/internal/repository"
)

func ServicesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := repository.GetCategories(db)
		if err != nil {
			http.Error(w, "Ошибка при получении категорий", http.StatusInternalServerError)
			return
		}

		data := model.PageData{
			Title:           "Услуги",
			Categories:      categories,
			CurrentCategory: "",
		}

		renderTemplate(w, "services.html", &data, false)
	}
}

func ServicesCategoryHandler(db *sql.DB, categorySlug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем категории для навигации
		categories, err := repository.GetCategories(db)
		if err != nil {
			http.Error(w, "Ошибка при получении категорий", http.StatusInternalServerError)
			return
		}

		// Получаем услуги для выбранной категории
		services, err := repository.GetServicesByCategory(db, categorySlug)
		if err != nil {
			http.Error(w, "Ошибка при получении услуг", http.StatusInternalServerError)
			return
		}

		// Находим текущую категорию
		var currentCategory model.Category
		for _, c := range categories {
			if c.Slug == categorySlug {
				currentCategory = c
				break
			}
		}

		data := model.PageData{
			Title:           currentCategory.Name,
			Categories:      categories,
			CurrentCategory: categorySlug,
			Services:        services,
		}

		renderTemplate(w, "services_category.html", &data, false)
	}
}

func ServiceDescriptionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем slug услуги из URL
		serviceSlug := r.URL.Path[len("/services/description/"):]

		// Получаем категории для навигации
		categories, err := repository.GetCategories(db)
		if err != nil {
			http.Error(w, "Ошибка при получении категорий", http.StatusInternalServerError)
			return
		}

		// Получаем информацию об услуге
		service, err := repository.GetServiceBySlug(db, serviceSlug)
		if err != nil {
			http.Error(w, "Услуга не найдена", http.StatusNotFound)
			return
		}

		data := model.PageData{
			Title:           service.Name,
			Categories:      categories,
			CurrentCategory: service.Category.Slug,
			Services:        []model.Service{service},
		}

		renderTemplate(w, "services/description/"+serviceSlug+".html", &data, false)
	}
}
