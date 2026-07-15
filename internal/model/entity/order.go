package entity

import "time"

type Order struct {
	ID        int
	OrderCode string
	UserID    int
	Total     float64
	Status    string
	CreatedAt time.Time
	Items     []OrderItem
}

type OrderItem struct {
	ID          int
	OrderID     int
	ServiceID   string
	ServiceName string
	Price       float64
}
