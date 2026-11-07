package dto

type CreateOrderRequest struct {
	Customer string            `json:"customer"`
	Items    []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	SKU string `json:"sku"`
	Qty int64  `json:"qty"`
}

type OrderResponse struct {
    ID        int64               `json:"id"`
    Customer  string              `json:"customer"`
    Status    string              `json:"status"`
    Total     int64               `json:"total"`
    Items     []OrderItemResponse `json:"items"`
    CreatedAt string              `json:"createdAt"`
}

type OrderItemResponse struct {
	SKU       string `json:"sku"`
	Qty       int64  `json:"qty"`
	UnitPrice int64  `json:"unitPrice"`
	Subtotal  int64  `json:"subtotal"`
}

type OrderListItemResponse struct {
    ID        int64  `json:"id"`
    Customer  string `json:"customer"`
    Status    string `json:"status"`
    Total     int64  `json:"total"`
    CreatedAt string `json:"createdAt"`
}

type OrderListResponse struct {
    Data  []OrderListItemResponse `json:"data"`
    Page  int                    `json:"page"`
    Limit int                    `json:"limit"`
    Total int64                  `json:"total"`
}