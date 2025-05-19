package repository

import (
	"database/sql"
	"myproject/internal/model"
)

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
