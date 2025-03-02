-- name: CreateRun :one
INSERT INTO runs (name, title, description, docker_image, parameters, data, mounts, created_at)
VALUES (?,?,?,?,?,?,?,datetime('now'))
RETURNING *;

-- name: GetRun :one
SELECT * FROM runs WHERE id = ?;

-- name: DeleteRun :exec
DELETE FROM runs WHERE id = ?;

-- name: StartRun :one
UPDATE runs SET status = 'running', started_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: FinishRun :one
UPDATE runs SET status = 'finished', finished_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: RunErrored :one
UPDATE runs SET status = 'errored', error_message = ?, finished_at = datetime('now'), has_errored = TRUE
WHERE id = ?
RETURNING *;

-- name: GetAllRuns :many
SELECT * FROM runs;

-- name: GetIdleRuns :many
SELECT * FROM runs WHERE status = 'pending';

-- name: GetRunning :many
SELECT * FROM runs WHERE status = 'running';

-- name: GetFinishedRuns :many
SELECT * FROM runs WHERE status = 'finished';

-- name: GetErroredRuns :many
SELECT * FROM runs WHERE status = 'errored';
