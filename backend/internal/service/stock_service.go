package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
)

type InsufficientStockError struct {
	Available int64
}

func (e *InsufficientStockError) Error() string {
	return "insufficient stock"
}

type StockMovementRepository interface {
	CreateMovement(
		ctx context.Context,
		productID int64,
		movement model.StockMovement,
	) (model.StockMovement, int64, error)

	GetMovements(
		ctx context.Context,
		productID int64,
		page int,
		limit int,
	) ([]model.StockMovement, int64, error)
}

type StockService struct {
	repo StockMovementRepository
}

func NewStockService(repo StockMovementRepository) *StockService {
	return &StockService{
		repo: repo,
	}
}

func (s *StockService) CreateMovement(
	ctx context.Context,
	productID int64,
	req dto.CreateStockMovementRequest,
) (dto.StockMovementResponse, error) {

	var response dto.StockMovementResponse

	if productID <= 0 {
		return response, ErrProductNotFound
	}

	req.Type = strings.ToLower(strings.TrimSpace(req.Type))

	if req.Type != "in" && req.Type != "out" {
		return response, ErrInvalidMovementType
	}

	if req.Qty <= 0 {
		return response, ErrInvalidQuantity
	}

	movement := model.StockMovement{
		ProductID: productID,
		Type:      req.Type,
		Qty:       req.Qty,
		Note:      strings.TrimSpace(req.Note),
	}
	result, available, err := s.repo.CreateMovement(
		ctx,
		productID,
		movement,
	)

	if err != nil {
		var insufficientErr *repository.InsufficientStockError

		if errors.As(err, &insufficientErr) {
			return response, &InsufficientStockError{
				Available: insufficientErr.Available,
			}
		}

		if errors.Is(err, repository.ErrProductNotFound) {
			return response, ErrProductNotFound
		}

		return response, err
	}

	response = dto.StockMovementResponse{
		ID:        result.ID,
		ProductID: result.ProductID,
		Type:      result.Type,
		Qty:       result.Qty,
		Note:      result.Note,
		Available: available,
		CreatedAt: result.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *StockService) GetMovements(
	ctx context.Context,
	productID int64,
	req dto.StockMovementListRequest,
) (dto.StockMovementListResponse, error) {

	response := dto.StockMovementListResponse{}

	if productID <= 0 {
		return response, ErrProductNotFound
	}

	movements, total, err := s.repo.GetMovements(
		ctx,
		productID,
		req.Page,
		req.Limit,
	)

	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return response, ErrProductNotFound
		}

		return response, err
	}

	data := make([]dto.StockMovementResponse, 0, len(movements))

	for _, movement := range movements {
		data = append(data, dto.StockMovementResponse{
			ID:        movement.ID,
			ProductID: movement.ProductID,
			Type:      movement.Type,
			Qty:       movement.Qty,
			Note:      movement.Note,
			CreatedAt: movement.CreatedAt.Format(time.RFC3339),
		})
	}

	response = dto.StockMovementListResponse{
		Data:  data,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return response, nil
}


