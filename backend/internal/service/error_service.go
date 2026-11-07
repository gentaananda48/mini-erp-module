package service

import "errors"

var (
	ErrCustomerRequired = errors.New("customer is required")
	ErrItemsRequired    = errors.New("order items are required")
	ErrDuplicateSKU     = errors.New("duplicate sku in order")
	ErrIdempotencyKeyRequired = errors.New("Idempotency-Key header is required")
	ErrSKURequired  = errors.New("sku is required")
	ErrSKUExists    = errors.New("sku already exists")
	ErrNameRequired = errors.New("name is required")
	ErrPriceRequired = errors.New("price is required")
	ErrInvalidPrice = errors.New("price must be an integer greater than or equal to 0")
	ErrProductNotFound = errors.New("product not found")
	ErrProductHasHistory = errors.New("product has movement history")
	ErrInvalidMovementType = errors.New("invalid movement type")
	ErrInvalidQuantity     = errors.New("quantity must be greater than zero")
    ErrInsufficientStock   = errors.New("insufficient stock")
    ErrOrderNotFound         = errors.New("order not found")
    ErrOrderAlreadyCancelled = errors.New("order already cancelled")
    ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled")	
)