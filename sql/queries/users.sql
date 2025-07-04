-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, hashed_password, email)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUserPwdEmail :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = $4
WHERE id = $1
RETURNING *;