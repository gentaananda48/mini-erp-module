package model

import "time"

type StockMovement struct {
	ID        int64
	ProductID int64
	Type      string
	Qty       int64
	Note      string
	CreatedAt time.Time
}