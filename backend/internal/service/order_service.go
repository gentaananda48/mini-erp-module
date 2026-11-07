package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	req dto.CreateOrderRequest,
	idempotencyKey string,
) (dto.OrderResponse, error) {

	var response dto.OrderResponse

	req.Customer = strings.TrimSpace(req.Customer)
	idempotencyKey = strings.TrimSpace(idempotencyKey)

	if req.Customer == "" {
		return response, ErrCustomerRequired
	}

	if len(req.Items) == 0 {
		return response, ErrItemsRequired
	}

	if idempotencyKey == "" {
		return response, ErrIdempotencyKeyRequired
	}

	seen := make(map[string]bool)

	order := model.Order{
		Customer: req.Customer,
		Items:    make([]model.OrderItem, 0, len(req.Items)),
	}

	for _, item := range req.Items {
		item.SKU = strings.TrimSpace(item.SKU)

		if item.SKU == "" {
			return response, ErrSKURequired
		}

		if item.Qty <= 0 {
			return response, ErrInvalidQuantity
		}

		if seen[item.SKU] {
			return response, ErrDuplicateSKU
		}

		seen[item.SKU] = true

		order.Items = append(order.Items, model.OrderItem{
			SKU: item.SKU,
			Qty: item.Qty,
		})
	}

	hashPayload := struct {
		Customer string
		Items    []model.OrderItem
	}{
		Customer: order.Customer,
		Items:    order.Items,
	}

	payload, err := json.Marshal(hashPayload)
	if err != nil {
		return response, err
	}

	sum := sha256.Sum256(payload)
	requestHash := hex.EncodeToString(sum[:])

	result, _, err := s.repo.CreateOrder(
		ctx,
		order,
		idempotencyKey,
		requestHash,
	)

	if err != nil {
		return response, err
	}

	response = toOrderResponse(result)

	return response, nil
}

func toOrderResponse(order model.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, dto.OrderItemResponse{
			SKU:       item.SKU,
			Qty:       item.Qty,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Qty * item.UnitPrice,
		})
	}

	return dto.OrderResponse{
		ID:        order.ID,
		Customer:  order.Customer,
		Status:    order.Status,
		Total:     order.Total,
		Items:     items,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}


func (s *OrderService) CancelOrder(
    ctx context.Context,
    orderID int64,
) error {

    if orderID <= 0 {
        return ErrOrderNotFound
    }

    err := s.repo.CancelOrder(ctx, orderID)

    if err != nil {
        switch {
        case errors.Is(err, repository.ErrOrderNotFound):
            return ErrOrderNotFound

        case errors.Is(err, repository.ErrOrderAlreadyCancelled):
            return ErrOrderAlreadyCancelled

        case errors.Is(err, repository.ErrOrderCannotBeCancelled):
            return ErrOrderCannotBeCancelled
        }

        return err
    }

    return nil
}

func (s *OrderService) GetOrderByID(
    ctx context.Context,
    orderID int64,
) (dto.OrderResponse, error) {
    if orderID <= 0 {
        return dto.OrderResponse{}, errors.New("invalid order ID")
    }

    order, err := s.repo.GetOrderByID(ctx, orderID)
    if err != nil {
        return dto.OrderResponse{}, err
    }

    return toOrderResponse(order), nil
}

func (s *OrderService) GetOrders(
    ctx context.Context,
    page int,
    limit int,
) (dto.OrderListResponse, error) {

    orders, total, err := s.repo.GetOrders(ctx, page, limit)
    if err != nil {
        return dto.OrderListResponse{}, err
    }

    data := make([]dto.OrderListItemResponse, 0, len(orders))

    for _, order := range orders {
        data = append(data, dto.OrderListItemResponse{
            ID:        order.ID,
            Customer:  order.Customer,
            Status:    order.Status,
            Total:     order.Total,
            CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        })
    }

    return dto.OrderListResponse{
        Data:  data,
        Page:  page,
        Limit: limit,
        Total: total,
    }, nil
}