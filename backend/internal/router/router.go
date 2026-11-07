package router

import (
	"net/http"

	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	stockRepo := repository.NewStockRepository(db)
    stockService := service.NewStockService(stockRepo)
    stockHandler := handler.NewStockHandler(stockService)

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

    auditRepo := repository.NewAuditRepository(db)

	// Products
	mux.HandleFunc(
		"GET /api/products",
		productHandler.GetProducts,
	)
	mux.HandleFunc(
		"POST /api/products",
		productHandler.CreateProduct,
	)
	mux.HandleFunc(
		"PUT /api/products/{id}", 
		productHandler.UpdateProduct,
	)
	mux.HandleFunc(
		"DELETE /api/products/{id}",
		productHandler.DeleteProduct,
	)

	// Stock
	mux.HandleFunc(
		"GET /api/products/{id}/movements",
		stockHandler.GetMovements,
	)
    mux.HandleFunc(
        "POST /api/products/{id}/stock",
        stockHandler.CreateMovement,
    )

	// Orders
	mux.HandleFunc(
		"POST /api/orders", 
		orderHandler.CreateOrder,
	)
	mux.HandleFunc(
		"POST /api/orders/{id}/cancel",
		orderHandler.CancelOrder,
	)
	mux.HandleFunc(
		"GET /api/orders/{id}", 
		orderHandler.GetOrderByID,
	)

    auditedHandler := middleware.Audit(auditRepo)(mux)

    return middleware.Logger(auditedHandler)
}
