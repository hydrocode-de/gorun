// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, is_admin)
VALUES (
    ?1,
    ?2,
    ?3,
    ?4
)
RETURNING id, email, password_hash, is_admin, created_at, last_login
`

type CreateUserParams struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	IsAdmin      bool   `json:"isAdmin"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.PasswordHash,
		arg.IsAdmin,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?1
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT users.id, users.email, users.is_admin, COALESCE(r.run_count, 0) as run_count
FROM users
LEFT JOIN (
    SELECT user_id, COUNT(id) as run_count
    FROM runs 
    GROUP BY user_id
) as r ON r.user_id = users.id
`

type GetAllUsersRow struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
	RunCount int64  `json:"runCount"`
}

func (q *Queries) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersRow
	for rows.Next() {
		var i GetAllUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.IsAdmin,
			&i.RunCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password_hash, is_admin, created_at, last_login FROM users
WHERE email = ?1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, password_hash, is_admin, created_at, last_login FROM users
WHERE id = ?1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users 
SET password_hash = ?1
WHERE id = ?2
RETURNING id, email, password_hash, is_admin, created_at, last_login
`

type UpdateUserPasswordParams struct {
	PasswordHash string `json:"passwordHash"`
	ID           string `json:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.PasswordHash, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.LastLogin,
	)
	return i, err
}
