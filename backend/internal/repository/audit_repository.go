package repository

import (
	"context"

	"backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
    db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
    return &AuditRepository{
        db: db,
    }
}

func (r *AuditRepository) Create(
    ctx context.Context,
    audit model.AuditLog,
) error {
    query := `
        INSERT INTO audit_logs (
            method,
            endpoint,
            resource_type,
            resource_id,
            request_body,
            response_status,
            response_body
        )
        VALUES (
            $1,
            $2,
            NULLIF($3, ''),
            NULLIF($4, ''),
            NULLIF($5, ''),
            $6,
            NULLIF($7, '')
        )
    `

    _, err := r.db.Exec(
        ctx,
        query,
        audit.Method,
        audit.Endpoint,
        audit.ResourceType,
        audit.ResourceID,
        audit.RequestBody,
        audit.ResponseStatus,
        audit.ResponseBody,
    )

    return err
}