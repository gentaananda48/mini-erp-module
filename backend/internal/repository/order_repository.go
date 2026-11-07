package repository

import (
	"context"
	"errors"
	"math"
	"sort"
	"strconv"

	"backend/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUnknownSKU             = errors.New("unknown sku")
	ErrIdempotencyKeyConflict = errors.New("idempotency key already used with different request")
    ErrOrderNotFound         = errors.New("order not found")
    ErrOrderAlreadyCancelled = errors.New("order already cancelled")
    ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled")
)

type UnknownSKUError struct {
	SKU string
}

func (e *UnknownSKUError) Error() string {
	return "unknown sku: " + e.SKU
}

type InsufficientOrderStockError struct {
	SKU       string
	Available int64
}

func (e *InsufficientOrderStockError) Error() string {
	return "insufficient stock"
}

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

type ProductForOrder struct {
	ID    int64
	SKU   string
	Price int64
	Stock int64
}

func (r *OrderRepository) CreateOrder(
	ctx context.Context,
	order model.Order,
	idempotencyKey string,
	requestHash string,
) (model.Order, bool, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.Order{}, false, err
	}
	defer tx.Rollback(ctx)

	var insertedKey string

	err = tx.QueryRow(ctx, `
		INSERT INTO order_idempotency (
			idempotency_key,
			request_hash
		)
		VALUES ($1, $2)
		ON CONFLICT (idempotency_key) DO NOTHING
		RETURNING idempotency_key
	`, idempotencyKey, requestHash).Scan(&insertedKey)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return model.Order{}, false, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		var existingHash string
		var existingOrderID int64

		err = tx.QueryRow(ctx, `
			SELECT request_hash, order_id
			FROM order_idempotency
			WHERE idempotency_key = $1
		`, idempotencyKey).Scan(&existingHash, &existingOrderID)

		if err != nil {
			return model.Order{}, false, err
		}

		if existingHash != requestHash {
			return model.Order{}, false, ErrIdempotencyKeyConflict
		}

		if existingOrderID == 0 {
			return model.Order{}, false, errors.New("idempotency record has no order")
		}

		existingOrder, err := getOrderByID(ctx, tx, existingOrderID)
		if err != nil {
			return model.Order{}, false, err
		}

		if err := tx.Commit(ctx); err != nil {
			return model.Order{}, false, err
		}

		return existingOrder, true, nil
	}

	skus := make([]string, 0, len(order.Items))
	for _, item := range order.Items {
		skus = append(skus, item.SKU)
	}
	sort.Strings(skus)

	rows, err := tx.Query(ctx, `
		SELECT id, sku, price, stock
		FROM products
		WHERE sku = ANY($1)
		ORDER BY sku
		FOR UPDATE
	`, skus)

	if err != nil {
		return model.Order{}, false, err
	}

	products := make(map[string]ProductForOrder)

	for rows.Next() {
		var product ProductForOrder

		if err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Price,
			&product.Stock,
		); err != nil {
			rows.Close()
			return model.Order{}, false, err
		}

		products[product.SKU] = product
	}

	rowsErr := rows.Err()
	rows.Close()

	if rowsErr != nil {
		return model.Order{}, false, rowsErr
	}

	for _, item := range order.Items {
		product, exists := products[item.SKU]

		if !exists {
			return model.Order{}, false, &UnknownSKUError{
				SKU: item.SKU,
			}
		}

		if product.Stock < item.Qty {
			return model.Order{}, false, &InsufficientOrderStockError{
				SKU:       item.SKU,
				Available: product.Stock,
			}
		}

		if product.Price < 0 || item.Qty > 0 && product.Price > math.MaxInt64/item.Qty {
			return model.Order{}, false, errors.New("order subtotal overflow")
		}

		subtotal := product.Price * item.Qty

		if order.Total > math.MaxInt64-subtotal {
			return model.Order{}, false, errors.New("order total overflow")
		}

		order.Total += subtotal
	}

	for _, item := range order.Items {
		product := products[item.SKU]

		tag, err := tx.Exec(ctx, `
			UPDATE products
			SET stock = stock - $1
			WHERE id = $2
			  AND stock >= $1
		`, item.Qty, product.ID)

		if err != nil {
			return model.Order{}, false, err
		}

		if tag.RowsAffected() != 1 {
			return model.Order{}, false, &InsufficientOrderStockError{
				SKU:       item.SKU,
				Available: product.Stock,
			}
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO stock_movements (
				product_id,
				type,
				qty,
				note
			)
			VALUES ($1, 'out', $2, $3)
		`, product.ID, item.Qty, "Sales order")

		if err != nil {
			return model.Order{}, false, err
		}
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO sales_orders (
			customer,
			total
		)
		VALUES ($1, $2)
		RETURNING id, created_at
	`, order.Customer, order.Total).Scan(
		&order.ID,
		&order.CreatedAt,
	)

	if err != nil {
		return model.Order{}, false, err
	}

	for i := range order.Items {
		item := &order.Items[i]
		product := products[item.SKU]

		item.OrderID = order.ID
		item.UnitPrice = product.Price

		err = tx.QueryRow(ctx, `
			INSERT INTO sales_order_items (
				order_id,
				sku,
				qty,
				unit_price
			)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, order.ID, item.SKU, item.Qty, item.UnitPrice).Scan(&item.ID)

		if err != nil {
			return model.Order{}, false, err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE order_idempotency
		SET order_id = $1
		WHERE idempotency_key = $2
	`, order.ID, idempotencyKey)

	if err != nil {
		return model.Order{}, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Order{}, false, err
	}

	return order, false, nil
}

func getOrderByID(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
) (model.Order, error) {

	var order model.Order

	err := tx.QueryRow(ctx, `
		SELECT id, customer, total, created_at
		FROM sales_orders
		WHERE id = $1
	`, orderID).Scan(
		&order.ID,
		&order.Customer,
		&order.Total,
		&order.CreatedAt,
	)

	if err != nil {
		return order, err
	}

	order.Items = make([]model.OrderItem, 0)

	rows, err := tx.Query(ctx, `
		SELECT id, order_id, sku, qty, unit_price
		FROM sales_order_items
		WHERE order_id = $1
		ORDER BY id
	`, orderID)

	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.SKU,
			&item.Qty,
			&item.UnitPrice,
		); err != nil {
			return order, err
		}

		order.Items = append(order.Items, item)
	}

	return order, rows.Err()
}


func (r *OrderRepository) CancelOrder(
    ctx context.Context,
    orderID int64,
) error {

    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    var status string

    err = tx.QueryRow(ctx, `
        SELECT status
        FROM sales_orders
        WHERE id = $1
        FOR UPDATE
    `, orderID).Scan(&status)

    if errors.Is(err, pgx.ErrNoRows) {
        return ErrOrderNotFound
    }

    if err != nil {
        return err
    }

    if status == "CANCELLED" {
        return ErrOrderAlreadyCancelled
    }

    if status != "CREATED" {
        return ErrOrderCannotBeCancelled
    }

    rows, err := tx.Query(ctx, `
        SELECT
            p.id,
            p.sku,
            soi.qty
        FROM sales_order_items soi
        JOIN products p ON p.sku = soi.sku
        WHERE soi.order_id = $1
        ORDER BY p.sku
        FOR UPDATE OF p
    `, orderID)

    if err != nil {
        return err
    }

    type cancellationItem struct {
        ProductID int64
        SKU       string
        Qty       int64
    }

    items := make([]cancellationItem, 0)

    for rows.Next() {
        var item cancellationItem

        if err := rows.Scan(
            &item.ProductID,
            &item.SKU,
            &item.Qty,
        ); err != nil {
            rows.Close()
            return err
        }

        items = append(items, item)
    }

    rowsErr := rows.Err()
    rows.Close()

    if rowsErr != nil {
        return rowsErr
    }

    for _, item := range items {
        _, err := tx.Exec(ctx, `
            UPDATE products
            SET stock = stock + $1
            WHERE id = $2
        `, item.Qty, item.ProductID)

        if err != nil {
            return err
        }

        _, err = tx.Exec(ctx, `
            INSERT INTO stock_movements (
                product_id,
                type,
                qty,
                note
            )
            VALUES ($1, 'in', $2, $3)
        `,
            item.ProductID,
            item.Qty,
            "Order cancellation #"+strconv.FormatInt(orderID, 10),
        )

        if err != nil {
            return err
        }
    }

    _, err = tx.Exec(ctx, `
        UPDATE sales_orders
        SET status = 'CANCELLED'
        WHERE id = $1
    `, orderID)

    if err != nil {
        return err
    }

    return tx.Commit(ctx)
}
func (r *OrderRepository) GetOrderByID(
    ctx context.Context,
    orderID int64,
) (model.Order, error) {
    var order model.Order

    err := r.db.QueryRow(ctx, `
        SELECT id, customer, status, total, created_at
        FROM sales_orders
        WHERE id = $1
    `, orderID).Scan(
        &order.ID,
        &order.Customer,
        &order.Status,
        &order.Total,
        &order.CreatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return order, ErrOrderNotFound
    }

    if err != nil {
        return order, err
    }

    order.Items = make([]model.OrderItem, 0)

    rows, err := r.db.Query(ctx, `
        SELECT sku, qty, unit_price
        FROM sales_order_items
        WHERE order_id = $1
        ORDER BY id ASC
    `, orderID)

    if err != nil {
        return order, err
    }
    defer rows.Close()

    for rows.Next() {
        var item model.OrderItem

        err := rows.Scan(
            &item.SKU,
            &item.Qty,
            &item.UnitPrice,
        )

        if err != nil {
            return order, err
        }

        order.Items = append(order.Items, item)
    }

    if err := rows.Err(); err != nil {
        return order, err
    }

    return order, nil
}

func (r *OrderRepository) GetOrders(
    ctx context.Context,
    page int,
    limit int,
) ([]model.Order, int64, error) {
    offset := (page - 1) * limit

    var total int64

    err := r.db.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM sales_orders
    `).Scan(&total)

    if err != nil {
        return nil, 0, err
    }

    rows, err := r.db.Query(ctx, `
        SELECT
            id,
            customer,
            status,
            total,
            created_at
        FROM sales_orders
        ORDER BY created_at DESC, id DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)

    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    orders := make([]model.Order, 0)

    for rows.Next() {
        var order model.Order

        if err := rows.Scan(
            &order.ID,
            &order.Customer,
            &order.Status,
            &order.Total,
            &order.CreatedAt,
        ); err != nil {
            return nil, 0, err
        }

        orders = append(orders, order)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, err
    }

    return orders, total, nil
}