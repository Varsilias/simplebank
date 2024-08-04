-- name: CreateEntry :one
INSERT INTO entries (
  public_id, account_id, amount, type, last_balance
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: GetEntryByAccountId :one
SELECT * FROM entries
WHERE account_id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListEntriesForAccountId :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id = $1;