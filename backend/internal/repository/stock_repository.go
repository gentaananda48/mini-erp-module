package repository

import (
	"context"
	"errors"
	"fmt"

	"backend/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
var ErrProductNotFound = errors.New("product not found")
type StockRepository struct {
	db *pgxpool.Pool
}

func NewStockRepository(db *pgxpool.Pool) *StockRepository {
	return &StockRepository{
		db: db,
	}
}

type InsufficientStockError struct {
	Available int64
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf(
		"insufficient stock: available %d",
		e.Available,
	)
}

func (r *StockRepository) CreateMovement(
	ctx context.Context,
	productID int64,
	movement model.StockMovement,
) (model.StockMovement, int64, error) {

	var result model.StockMovement

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return result, 0, err
	}
	defer tx.Rollback(ctx)

	var available int64

	err = tx.QueryRow(ctx, `
		SELECT stock
		FROM products
		WHERE id = $1
		FOR UPDATE
	`, productID).Scan(&available)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, 0, ErrProductNotFound
		}

		return result, 0, err
	}

	var newStock int64

	switch movement.Type {
	case "in":
		newStock = available + movement.Qty

	case "out":
		if movement.Qty > available {
			return result, available, &InsufficientStockError{
				Available: available,
			}
		}

		newStock = available - movement.Qty

	default:
		return result, available, errors.New("invalid movement type")
	}

	_, err = tx.Exec(ctx, `
		UPDATE products
		SET stock = $1
		WHERE id = $2
	`, newStock, productID)

	if err != nil {
		return result, available, err
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO stock_movements (
			product_id,
			type,
			qty,
			note
		)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		RETURNING id, created_at
	`,
		productID,
		movement.Type,
		movement.Qty,
		movement.Note,
	).Scan(
		&result.ID,
		&result.CreatedAt,
	)

	if err != nil {
		return result, available, err
	}

	result.ProductID = productID
	result.Type = movement.Type
	result.Qty = movement.Qty
	result.Note = movement.Note

	if err := tx.Commit(ctx); err != nil {
		return model.StockMovement{}, 0, err
	}

	return result, newStock, nil
}

func (r *StockRepository) GetMovements(
	ctx context.Context,
	productID int64,
	page int,
	limit int,
) ([]model.StockMovement, int64, error) {

	var exists bool

	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM products
			WHERE id = $1
		)
	`, productID).Scan(&exists)

	if err != nil {
		return nil, 0, err
	}

	if !exists {
		return nil, 0, ErrProductNotFound
	}

	var total int64

	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM stock_movements
		WHERE product_id = $1
	`, productID).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			product_id,
			type,
			qty,
			COALESCE(note, ''),
			created_at
		FROM stock_movements
		WHERE product_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`, productID, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movements := make([]model.StockMovement, 0)

	for rows.Next() {
		var movement model.StockMovement

		err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.Type,
			&movement.Qty,
			&movement.Note,
			&movement.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		movements = append(movements, movement)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}