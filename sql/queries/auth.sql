-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, created_at, expires_at)
VALUES (
    @user_id,
    @token,
    datetime('now'),
    @expires_at
)
RETURNING *;


-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET is_revoked = TRUE
WHERE token = @token;

-- name: GetUserRefreshTokens :many
SELECT * FROM refresh_tokens
WHERE user_id = @user_id AND is_revoked = FALSE;

-- name: GetRefreshTokenUser :one
SELECT u.* FROM users u
JOIN refresh_tokens rt ON rt.user_id = u.id
WHERE rt.token = @token
AND rt.expires_at > datetime('now')
AND rt.is_revoked = FALSE; 