-- name: CreateRun :one
INSERT INTO runs (name, title, description, docker_image, parameters, data, mounts, created_at)
VALUES (?,?,?,?,?,?,?,date('now'))
RETURNING *;

-- name: GetRun :one
SELECT * FROM runs WHERE id = ?;

-- name: StartRun :one
UPDATE runs SET status = 'running', started_at = date('now')
WHERE id = ?
RETURNING *;

-- name: FinishRun :one
UPDATE runs SET status = 'finished', finished_at = date('now')
WHERE id = ?
RETURNING *;

-- name: RunErrored :one
UPDATE runs SET status = 'errored', error_message = ?, finished_at = date('now'), has_errored = TRUE
WHERE id = ?
RETURNING *;

-- name: GetAllRuns :many
SELECT * FROM runs;

-- name: GetIdleRuns :many
SELECT * FROM runs WHERE is_pending = TRUE;

-- name: GetRunning :many
SELECT * FROM runs WHERE is_running = TRUE;
