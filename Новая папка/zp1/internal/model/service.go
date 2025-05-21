package model

// Category представляет категорию услуг
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Service представляет услугу в системе
type Service struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	Slug       string   `json:"slug"`
	CategoryID int      `json:"category_id"`
	Category   Category `json:"category"`
}
