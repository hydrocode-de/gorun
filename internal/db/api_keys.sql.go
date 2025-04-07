// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: api_keys.sql

package db

import (
	"context"
	"database/sql"
)

const createApiKey = `-- name: CreateApiKey :one
INSERT INTO api_keys (key, created_at, valid_until)
VALUES (
    ?1,
    datetime('now'),
    ?2
)
RETURNING id, "key", created_at, last_used, valid_until
`

type CreateApiKeyParams struct {
	Key        string       `json:"key"`
	ValidUntil sql.NullTime `json:"validUntil"`
}

func (q *Queries) CreateApiKey(ctx context.Context, arg CreateApiKeyParams) (ApiKey, error) {
	row := q.db.QueryRowContext(ctx, createApiKey, arg.Key, arg.ValidUntil)
	var i ApiKey
	err := row.Scan(
		&i.ID,
		&i.Key,
		&i.CreatedAt,
		&i.LastUsed,
		&i.ValidUntil,
	)
	return i, err
}

const updateApiKeyLastUsed = `-- name: UpdateApiKeyLastUsed :exec
UPDATE api_keys
SET last_used = datetime('now')
WHERE key = ?1 AND valid_until > datetime('now')
`

func (q *Queries) UpdateApiKeyLastUsed(ctx context.Context, key string) error {
	_, err := q.db.ExecContext(ctx, updateApiKeyLastUsed, key)
	return err
}

const validateApiKey = `-- name: ValidateApiKey :one
SELECT
    CASE
        WHEN EXISTS (
            SELECT 1 FROM api_keys
            WHERE api_keys.key = ?1 AND api_keys.valid_until > datetime('now')
        ) THEN 'valid'
        WHEN EXISTS (
            SELECT 1 FROM api_keys
            WHERE api_keys.key = ?1 AND api_keys.valid_until <= datetime('now')
        ) THEN 'expired'
        ELSE 'invalid'
    END as status
`

func (q *Queries) ValidateApiKey(ctx context.Context, key string) (string, error) {
	row := q.db.QueryRowContext(ctx, validateApiKey, key)
	var status string
	err := row.Scan(&status)
	return status, err
}
