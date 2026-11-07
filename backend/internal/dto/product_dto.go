package dto

type ProductListRequest struct {
	Search string
	Page   int
	Limit  int
}

type CreateProductRequest struct {
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price *int64 `json:"price"`
}

type UpdateProductRequest struct {
    SKU   string `json:"sku"`
    Name  string `json:"name"`
    Price *int64 `json:"price"`
}

type ProductResponse struct {
	ID             int64      `json:"id"`
	SKU            string     `json:"sku"`
	Name           string     `json:"name"`
	Price          int64      `json:"price"`
	Stock          int64      `json:"stock"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      string     `json:"created_at"`
	LastMovementAt *string    `json:"last_movement_at"`
}

type ProductListResponse struct {
	Data  []ProductResponse `json:"data"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
	Total int64             `json:"total"`
}
