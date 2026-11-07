package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func toProductResponse(p model.Product) dto.ProductResponse {
	response := dto.ProductResponse{
		ID:        p.ID,
		SKU:       p.SKU,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		IsActive:  p.IsActive,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	}

	if p.LastMovementAt != nil {
		value := p.LastMovementAt.Format(time.RFC3339)
		response.LastMovementAt = &value
	}

	return response
}

func (s *ProductService) GetProducts(
	ctx context.Context,
	req dto.ProductListRequest,
) (*dto.ProductListResponse, error) {

	if req.Page < 1 {
		return nil, errors.New("page must be greater than zero")
	}

	if req.Limit < 1 || req.Limit > 100 {
		return nil, errors.New("limit must be between 1 and 100")
	}

	req.Search = strings.TrimSpace(req.Search)

	offset := (req.Page - 1) * req.Limit

	products, total, err := s.repo.FindAll(
		ctx,
		req.Search,
		req.Limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	data := make([]dto.ProductResponse, 0, len(products))

	for _, p := range products {
		data = append(data, toProductResponse(p))
	}

	return &dto.ProductListResponse{
		Data:  data,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}, nil
}

func (s *ProductService) CreateProduct(
	ctx context.Context,
	req dto.CreateProductRequest,
) (*dto.ProductResponse, error) {

	req.SKU = strings.TrimSpace(req.SKU)
	req.Name = strings.TrimSpace(req.Name)

	if req.SKU == "" {
		return nil, ErrSKURequired
	}

	if req.Name == "" {
		return nil, ErrNameRequired
	}

	if req.Price == nil {
		return nil, ErrPriceRequired
	}

	if *req.Price < 0 {
		return nil, ErrInvalidPrice
	}

	product, err := s.repo.Create(
		ctx,
		req.SKU,
		req.Name,
		*req.Price,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "products_sku_key" {
			return nil, ErrSKUExists
		}

		return nil, err
	}

	response := toProductResponse(product)

	return &response, nil
}

func (s *ProductService) UpdateProduct(
    ctx context.Context,
    id int64,
    req dto.UpdateProductRequest,
) (*dto.ProductResponse, error) {

    req.SKU = strings.TrimSpace(req.SKU)
    req.Name = strings.TrimSpace(req.Name)

    if req.SKU == "" {
        return nil, ErrSKURequired
    }

    if req.Name == "" {
        return nil, ErrNameRequired
    }

    if req.Price == nil {
        return nil, ErrPriceRequired
    }

    if *req.Price < 0 {
        return nil, ErrInvalidPrice
    }

    product, err := s.repo.Update(
        ctx,
        id,
        req.SKU,
        req.Name,
        *req.Price,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrProductNotFound
    }

    if err != nil {
        return nil, err
    }

    response := toProductResponse(product)

    return &response, nil
}

func (s *ProductService) DeleteProduct(
    ctx context.Context,
    id int64,
) error {
    exists, err := s.repo.HasHistory(ctx, id)
    if err != nil {
        return err
    }

    if exists {
        return ErrProductHasHistory
    }

    return s.repo.Delete(ctx, id)
}