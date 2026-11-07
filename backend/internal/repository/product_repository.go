package repository

import (
	"context"

	"backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll(
	ctx context.Context,
	search string,
	limit int,
	offset int,
) ([]model.Product, int64, error) {

	var total int64

	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM products p
		WHERE (
			$1 = ''
			OR p.sku ILIKE '%' || $1 || '%'
			OR p.name ILIKE '%' || $1 || '%'
		)
	`, search).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			p.id,
			p.sku,
			p.name,
			p.price,
			p.stock,
			p.is_active,
			p.created_at,
			lm.last_movement_at
		FROM products p
		LEFT JOIN (
			SELECT
				product_id,
				MAX(created_at) AS last_movement_at
			FROM stock_movements
			GROUP BY product_id
		) lm ON lm.product_id = p.id
		WHERE (
			$1 = ''
			OR p.sku ILIKE '%' || $1 || '%'
			OR p.name ILIKE '%' || $1 || '%'
		)
		ORDER BY p.id DESC
		LIMIT $2 OFFSET $3
	`, search, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)

	for rows.Next() {
		var p model.Product

		err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.IsActive,
			&p.CreatedAt,
			&p.LastMovementAt,
		)

		if err != nil {
			return nil, 0, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}


func (r *ProductRepository) Create(
	ctx context.Context,
	sku string,
	name string,
	price int64,
) (model.Product, error) {

	var product model.Product

	query := `
		INSERT INTO products (
			sku,
			name,
			price
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			sku,
			name,
			price,
			stock,
			is_active,
			created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		sku,
		name,
		price,
	).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.IsActive,
		&product.CreatedAt,
	)

	if err != nil {
		return model.Product{}, err
	}

	product.LastMovementAt = nil

	return product, nil
}

func (r *ProductRepository) Update(
    ctx context.Context,
    id int64,
    sku string,
    name string,
    price int64,
) (model.Product, error) {
    var product model.Product

    query := `
        UPDATE products
        SET
            sku = $1,
            name = $2,
            price = $3
        WHERE id = $4
        RETURNING
            id,
            sku,
            name,
            price,
            stock,
            is_active,
            created_at
    `

    err := r.db.QueryRow(
        ctx,
        query,
        sku,
        name,
        price,
        id,
    ).Scan(
        &product.ID,
        &product.SKU,
        &product.Name,
        &product.Price,
        &product.Stock,
        &product.IsActive,
        &product.CreatedAt,
    )

    if err != nil {
        return model.Product{}, err
    }

    product.LastMovementAt = nil

    return product, nil
}

func (r *ProductRepository) HasHistory(
    ctx context.Context,
    id int64,
) (bool, error) {
    var exists bool

    err := r.db.QueryRow(
        ctx,
        `SELECT EXISTS (
            SELECT 1
            FROM stock_movements
            WHERE product_id = $1
        )`,
        id,
    ).Scan(&exists)

    return exists, err
}

func (r *ProductRepository) Delete(
    ctx context.Context,
    id int64,
) error {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    var productID int64

    err = tx.QueryRow(ctx, `
        SELECT id
        FROM products
        WHERE id = $1
        FOR UPDATE
    `, id).Scan(&productID)

    if err != nil {
        return err
    }

    _, err = tx.Exec(ctx, `
        DELETE FROM products
        WHERE id = $1
    `, id)

    if err != nil {
        return err
    }

    return tx.Commit(ctx)
}