-- name: CreateApiKey :one
INSERT INTO api_keys (key, created_at, valid_until)
VALUES (
    @key,
    datetime('now'),
    @valid_until
)
RETURNING *;

-- name: ValidateApiKey :one
SELECT
    CASE
        WHEN EXISTS (
            SELECT 1 FROM api_keys
            WHERE api_keys.key = @key AND api_keys.valid_until > datetime('now')
        ) THEN 'valid'
        WHEN EXISTS (
            SELECT 1 FROM api_keys
            WHERE api_keys.key = @key AND api_keys.valid_until <= datetime('now')
        ) THEN 'expired'
        ELSE 'invalid'
    END as status;

-- name: UpdateApiKeyLastUsed :exec
UPDATE api_keys
SET last_used = datetime('now')
WHERE key = @key AND valid_until > datetime('now');

-- name: GenerateLetterOnlyKey :one
SELECT substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) ||
       substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) ||
       substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) ||
       substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) ||
       substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) ||
       substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) ||
       substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) ||
       substr('ABCDEFGHIJKLMNOPQRSTUVWXYZ', abs(random()) % 26 + 1, 1) as key;