-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, expires_at, user_id)
    VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET updated_at = $1,
    expires_at = $2
WHERE user_id = $3
RETURNING *;
