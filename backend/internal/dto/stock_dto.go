package dto

type CreateStockMovementRequest struct {
	Type string `json:"type"`
	Qty  int64  `json:"qty"`
	Note string `json:"note"`
}

type StockMovementResponse struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"productId"`
	Type      string `json:"type"`
	Qty       int64  `json:"qty"`
	Note      string `json:"note"`
	Available int64  `json:"available"`
	CreatedAt string `json:"createdAt"`
}

type StockMovementListRequest struct {
	Page  int
	Limit int
}

type StockMovementListResponse struct {
	Data  []StockMovementResponse `json:"data"`
	Page  int                     `json:"page"`
	Limit int                     `json:"limit"`
	Total int64                   `json:"total"`
}