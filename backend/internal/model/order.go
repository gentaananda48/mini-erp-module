package model

import "time"

type Order struct {
    ID        int64
    Customer  string
    Status    string
    Total     int64
    CreatedAt time.Time
    Items     []OrderItem
}

type OrderItem struct {
	ID        int64
	OrderID   int64
	SKU       string
	Qty       int64
	UnitPrice int64
}
