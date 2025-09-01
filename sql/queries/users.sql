-- name: CreateUser :one
INSERT INTO users(id, username, hashed_password, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: DeleteUsers :exec
DELETE FROM users;
