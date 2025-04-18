// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: runs.sql

package db

import (
	"context"
	"database/sql"
)

const createRun = `-- name: CreateRun :one
INSERT INTO runs (name, title, description, docker_image, parameters, data, mounts, created_at, user_id)
VALUES (?,?,?,?,?,?,?,datetime('now'),?)
RETURNING id, name, title, description, docker_image, mounts, parameters, data, created_at, started_at, finished_at, status, has_errored, error_message, user_id
`

type CreateRunParams struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DockerImage string `json:"dockerImage"`
	Parameters  string `json:"parameters"`
	Data        string `json:"data"`
	Mounts      string `json:"mounts"`
	UserID      string `json:"userId"`
}

func (q *Queries) CreateRun(ctx context.Context, arg CreateRunParams) (Run, error) {
	row := q.db.QueryRowContext(ctx, createRun,
		arg.Name,
		arg.Title,
		arg.Description,
		arg.DockerImage,
		arg.Parameters,
		arg.Data,
		arg.Mounts,
		arg.UserID,
	)
	var i Run
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Title,
		&i.Description,
		&i.DockerImage,
		&i.Mounts,
		&i.Parameters,
		&i.Data,
		&i.CreatedAt,
		&i.StartedAt,
		&i.FinishedAt,
		&i.Status,
		&i.HasErrored,
		&i.ErrorMessage,
		&i.UserID,
	)
	return i, err
}

const deleteRun = `-- name: DeleteRun :exec
DELETE FROM runs
WHERE runs.id = ? AND (
  (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
  OR runs.user_id = ?
)
`

type DeleteRunParams struct {
	ID     int64  `json:"id"`
	ID_2   string `json:"id2"`
	UserID string `json:"userId"`
}

func (q *Queries) DeleteRun(ctx context.Context, arg DeleteRunParams) error {
	_, err := q.db.ExecContext(ctx, deleteRun, arg.ID, arg.ID_2, arg.UserID)
	return err
}

const finishRun = `-- name: FinishRun :one
UPDATE runs SET status = 'finished', finished_at = datetime('now')
WHERE runs.id = ?
RETURNING id, name, title, description, docker_image, mounts, parameters, data, created_at, started_at, finished_at, status, has_errored, error_message, user_id
`

func (q *Queries) FinishRun(ctx context.Context, id int64) (Run, error) {
	row := q.db.QueryRowContext(ctx, finishRun, id)
	var i Run
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Title,
		&i.Description,
		&i.DockerImage,
		&i.Mounts,
		&i.Parameters,
		&i.Data,
		&i.CreatedAt,
		&i.StartedAt,
		&i.FinishedAt,
		&i.Status,
		&i.HasErrored,
		&i.ErrorMessage,
		&i.UserID,
	)
	return i, err
}

const getAllRuns = `-- name: GetAllRuns :many
SELECT r.id, r.name, r.title, r.description, r.docker_image, r.mounts, r.parameters, r.data, r.created_at, r.started_at, r.finished_at, r.status, r.has_errored, r.error_message, r.user_id FROM runs r
WHERE (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
   OR r.user_id = ?
`

type GetAllRunsParams struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (q *Queries) GetAllRuns(ctx context.Context, arg GetAllRunsParams) ([]Run, error) {
	rows, err := q.db.QueryContext(ctx, getAllRuns, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Run
	for rows.Next() {
		var i Run
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Title,
			&i.Description,
			&i.DockerImage,
			&i.Mounts,
			&i.Parameters,
			&i.Data,
			&i.CreatedAt,
			&i.StartedAt,
			&i.FinishedAt,
			&i.Status,
			&i.HasErrored,
			&i.ErrorMessage,
			&i.UserID,
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

const getErroredRuns = `-- name: GetErroredRuns :many
SELECT r.id, r.name, r.title, r.description, r.docker_image, r.mounts, r.parameters, r.data, r.created_at, r.started_at, r.finished_at, r.status, r.has_errored, r.error_message, r.user_id FROM runs r
WHERE r.status = 'errored' AND (
  (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
  OR r.user_id = ?
)
`

type GetErroredRunsParams struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (q *Queries) GetErroredRuns(ctx context.Context, arg GetErroredRunsParams) ([]Run, error) {
	rows, err := q.db.QueryContext(ctx, getErroredRuns, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Run
	for rows.Next() {
		var i Run
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Title,
			&i.Description,
			&i.DockerImage,
			&i.Mounts,
			&i.Parameters,
			&i.Data,
			&i.CreatedAt,
			&i.StartedAt,
			&i.FinishedAt,
			&i.Status,
			&i.HasErrored,
			&i.ErrorMessage,
			&i.UserID,
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

const getFinishedRuns = `-- name: GetFinishedRuns :many
SELECT r.id, r.name, r.title, r.description, r.docker_image, r.mounts, r.parameters, r.data, r.created_at, r.started_at, r.finished_at, r.status, r.has_errored, r.error_message, r.user_id FROM runs r
WHERE r.status = 'finished' AND (
  (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
  OR r.user_id = ?
)
`

type GetFinishedRunsParams struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (q *Queries) GetFinishedRuns(ctx context.Context, arg GetFinishedRunsParams) ([]Run, error) {
	rows, err := q.db.QueryContext(ctx, getFinishedRuns, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Run
	for rows.Next() {
		var i Run
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Title,
			&i.Description,
			&i.DockerImage,
			&i.Mounts,
			&i.Parameters,
			&i.Data,
			&i.CreatedAt,
			&i.StartedAt,
			&i.FinishedAt,
			&i.Status,
			&i.HasErrored,
			&i.ErrorMessage,
			&i.UserID,
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

const getIdleRuns = `-- name: GetIdleRuns :many
SELECT r.id, r.name, r.title, r.description, r.docker_image, r.mounts, r.parameters, r.data, r.created_at, r.started_at, r.finished_at, r.status, r.has_errored, r.error_message, r.user_id FROM runs r
WHERE r.status = 'pending' AND (
  (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
  OR r.user_id = ?
)
`

type GetIdleRunsParams struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (q *Queries) GetIdleRuns(ctx context.Context, arg GetIdleRunsParams) ([]Run, error) {
	rows, err := q.db.QueryContext(ctx, getIdleRuns, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Run
	for rows.Next() {
		var i Run
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Title,
			&i.Description,
			&i.DockerImage,
			&i.Mounts,
			&i.Parameters,
			&i.Data,
			&i.CreatedAt,
			&i.StartedAt,
			&i.FinishedAt,
			&i.Status,
			&i.HasErrored,
			&i.ErrorMessage,
			&i.UserID,
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

const getRun = `-- name: GetRun :one
SELECT r.id, r.name, r.title, r.description, r.docker_image, r.mounts, r.parameters, r.data, r.created_at, r.started_at, r.finished_at, r.status, r.has_errored, r.error_message, r.user_id FROM runs r
WHERE r.id = ? AND (
  (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
  OR r.user_id = ?
)
`

type GetRunParams struct {
	ID     int64  `json:"id"`
	ID_2   string `json:"id2"`
	UserID string `json:"userId"`
}

func (q *Queries) GetRun(ctx context.Context, arg GetRunParams) (Run, error) {
	row := q.db.QueryRowContext(ctx, getRun, arg.ID, arg.ID_2, arg.UserID)
	var i Run
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Title,
		&i.Description,
		&i.DockerImage,
		&i.Mounts,
		&i.Parameters,
		&i.Data,
		&i.CreatedAt,
		&i.StartedAt,
		&i.FinishedAt,
		&i.Status,
		&i.HasErrored,
		&i.ErrorMessage,
		&i.UserID,
	)
	return i, err
}

const getRunning = `-- name: GetRunning :many
SELECT r.id, r.name, r.title, r.description, r.docker_image, r.mounts, r.parameters, r.data, r.created_at, r.started_at, r.finished_at, r.status, r.has_errored, r.error_message, r.user_id FROM runs r
WHERE r.status = 'running' AND (
  (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
  OR r.user_id = ?
)
`

type GetRunningParams struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

func (q *Queries) GetRunning(ctx context.Context, arg GetRunningParams) ([]Run, error) {
	rows, err := q.db.QueryContext(ctx, getRunning, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Run
	for rows.Next() {
		var i Run
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Title,
			&i.Description,
			&i.DockerImage,
			&i.Mounts,
			&i.Parameters,
			&i.Data,
			&i.CreatedAt,
			&i.StartedAt,
			&i.FinishedAt,
			&i.Status,
			&i.HasErrored,
			&i.ErrorMessage,
			&i.UserID,
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

const runErrored = `-- name: RunErrored :one
UPDATE runs SET status = 'errored', error_message = ?, finished_at = datetime('now'), has_errored = TRUE
WHERE runs.id = ?
RETURNING id, name, title, description, docker_image, mounts, parameters, data, created_at, started_at, finished_at, status, has_errored, error_message, user_id
`

type RunErroredParams struct {
	ErrorMessage sql.NullString `json:"errorMessage"`
	ID           int64          `json:"id"`
}

func (q *Queries) RunErrored(ctx context.Context, arg RunErroredParams) (Run, error) {
	row := q.db.QueryRowContext(ctx, runErrored, arg.ErrorMessage, arg.ID)
	var i Run
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Title,
		&i.Description,
		&i.DockerImage,
		&i.Mounts,
		&i.Parameters,
		&i.Data,
		&i.CreatedAt,
		&i.StartedAt,
		&i.FinishedAt,
		&i.Status,
		&i.HasErrored,
		&i.ErrorMessage,
		&i.UserID,
	)
	return i, err
}

const startRun = `-- name: StartRun :one
UPDATE runs
SET status = 'running', started_at = datetime('now')
WHERE runs.id = ? AND (
  (SELECT u.is_admin FROM users u WHERE u.id = ?) = TRUE 
  OR runs.user_id = ?
)
RETURNING id, name, title, description, docker_image, mounts, parameters, data, created_at, started_at, finished_at, status, has_errored, error_message, user_id
`

type StartRunParams struct {
	ID     int64  `json:"id"`
	ID_2   string `json:"id2"`
	UserID string `json:"userId"`
}

func (q *Queries) StartRun(ctx context.Context, arg StartRunParams) (Run, error) {
	row := q.db.QueryRowContext(ctx, startRun, arg.ID, arg.ID_2, arg.UserID)
	var i Run
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Title,
		&i.Description,
		&i.DockerImage,
		&i.Mounts,
		&i.Parameters,
		&i.Data,
		&i.CreatedAt,
		&i.StartedAt,
		&i.FinishedAt,
		&i.Status,
		&i.HasErrored,
		&i.ErrorMessage,
		&i.UserID,
	)
	return i, err
}
