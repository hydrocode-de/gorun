-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, is_admin)
VALUES (
    @id,
    @email,
    @password_hash,
    @is_admin
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = @email;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @id;

-- name: GetAllUsers :many
SELECT users.id, users.email, users.is_admin, COALESCE(r.run_count, 0) as run_count
FROM users
LEFT JOIN (
    SELECT user_id, COUNT(id) as run_count
    FROM runs 
    GROUP BY user_id
) as r ON r.user_id = users.id;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id;

-- name: UpdateUserPassword :one
UPDATE users 
SET password_hash = @password_hash
WHERE id = @id
RETURNING *;