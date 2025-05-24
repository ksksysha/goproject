package model

type SortData struct {
	Field     string
	Direction string
}

type PaginationData struct {
	CurrentPage  int
	TotalPages   int
	PageSize     int
	TotalRecords int
	HasNext      bool
	HasPrev      bool
}

type PageData struct {
	Title             string
	Username          string
	Role              string
	UserID            int
	Services          []Service
	Users             []User
	Bookings          []Booking
	Content           string
	SelectedServiceID string
	Pagination        PaginationData
	Sort              SortData
}
