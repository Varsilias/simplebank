-- name: CreateUser :one
INSERT INTO users (
  public_id, firstname, lastname, email, password, salt
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = sqlc.arg(email) LIMIT 1;

-- name: GetUserByPublicID :one
SELECT * FROM users
WHERE public_id = sqlc.arg(public_id) LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;