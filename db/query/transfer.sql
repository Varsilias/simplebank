-- name: CreateTransfer :one
INSERT INTO transfers (
  public_id, from_account_id, to_account_id, amount
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: GetTransferByFromAccountId :one
SELECT * FROM transfers
WHERE from_account_id = $1 LIMIT 1;

-- name: GetTransferByToAccountId :one
SELECT * FROM transfers
WHERE to_account_id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListTransfersForFromAccountId :many
SELECT * FROM transfers
WHERE from_account_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListTransfersForToAccountId :many
SELECT * FROM transfers
WHERE to_account_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id = $1;