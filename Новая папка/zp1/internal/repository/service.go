package repository

import (
	"database/sql"
	"myproject/internal/model"
)

func GetCategories(db *sql.DB) ([]model.Category, error) {
	rows, err := db.Query("SELECT id, name, slug FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func GetServicesByCategory(db *sql.DB, categorySlug string) ([]model.Service, error) {
	query := `
		SELECT s.id, s.name, s.price, s.category_id, c.name as category_name, c.slug as category_slug
		FROM services s
		JOIN categories c ON s.category_id = c.id
		WHERE c.slug = $1
		ORDER BY s.name`

	rows, err := db.Query(query, categorySlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.CategoryID, &s.Category.Name, &s.Category.Slug); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func GetAllServices(db *sql.DB) ([]model.Service, error) {
	query := `
		SELECT s.id, s.name, s.price, s.category_id, c.name as category_name, c.slug as category_slug
		FROM services s
		JOIN categories c ON s.category_id = c.id
		ORDER BY c.name, s.name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.CategoryID, &s.Category.Name, &s.Category.Slug); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func GetServices(db *sql.DB) ([]model.Service, error) {
	rows, err := db.Query("SELECT id, name, price FROM services")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Price); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func GetServiceBySlug(db *sql.DB, slug string) (model.Service, error) {
	query := `
		SELECT s.id, s.name, s.price, s.category_id, c.name as category_name, c.slug as category_slug
		FROM services s
		JOIN categories c ON s.category_id = c.id
		WHERE s.slug = $1`

	var service model.Service
	err := db.QueryRow(query, slug).Scan(
		&service.ID,
		&service.Name,
		&service.Price,
		&service.CategoryID,
		&service.Category.Name,
		&service.Category.Slug,
	)
	if err != nil {
		return model.Service{}, err
	}
	return service, nil
}

func GetServiceByID(db *sql.DB, id int) (model.Service, error) {
	query := `
		SELECT s.id, s.name, s.price, s.slug, s.category_id, c.name as category_name, c.slug as category_slug
		FROM services s
		JOIN categories c ON s.category_id = c.id
		WHERE s.id = $1`

	var service model.Service
	err := db.QueryRow(query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Price,
		&service.Slug,
		&service.CategoryID,
		&service.Category.Name,
		&service.Category.Slug,
	)
	if err != nil {
		return model.Service{}, err
	}
	return service, nil
}
