-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;