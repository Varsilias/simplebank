// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: transfer.sql

package db

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
  public_id, from_account_id, to_account_id, amount
) VALUES (
  $1, $2, $3, $4
) RETURNING id, from_account_id, to_account_id, public_id, created_at, updated_at, deleted_at, amount
`

type CreateTransferParams struct {
	PublicID      string `json:"public_id"`
	FromAccountID int32  `json:"from_account_id"`
	ToAccountID   int32  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer,
		arg.PublicID,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Amount,
	)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.PublicID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Amount,
	)
	return i, err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id = $1
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, public_id, created_at, updated_at, deleted_at, amount FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int32) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.PublicID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Amount,
	)
	return i, err
}

const getTransferByFromAccountId = `-- name: GetTransferByFromAccountId :one
SELECT id, from_account_id, to_account_id, public_id, created_at, updated_at, deleted_at, amount FROM transfers
WHERE from_account_id = $1 LIMIT 1
`

func (q *Queries) GetTransferByFromAccountId(ctx context.Context, fromAccountID int32) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransferByFromAccountId, fromAccountID)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.PublicID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Amount,
	)
	return i, err
}

const getTransferByToAccountId = `-- name: GetTransferByToAccountId :one
SELECT id, from_account_id, to_account_id, public_id, created_at, updated_at, deleted_at, amount FROM transfers
WHERE to_account_id = $1 LIMIT 1
`

func (q *Queries) GetTransferByToAccountId(ctx context.Context, toAccountID int32) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransferByToAccountId, toAccountID)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.PublicID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Amount,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, public_id, created_at, updated_at, deleted_at, amount FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListTransfersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.PublicID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Amount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfersForFromAccountId = `-- name: ListTransfersForFromAccountId :many
SELECT id, from_account_id, to_account_id, public_id, created_at, updated_at, deleted_at, amount FROM transfers
WHERE from_account_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3
`

type ListTransfersForFromAccountIdParams struct {
	FromAccountID int32 `json:"from_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *Queries) ListTransfersForFromAccountId(ctx context.Context, arg ListTransfersForFromAccountIdParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfersForFromAccountId, arg.FromAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.PublicID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Amount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfersForToAccountId = `-- name: ListTransfersForToAccountId :many
SELECT id, from_account_id, to_account_id, public_id, created_at, updated_at, deleted_at, amount FROM transfers
WHERE to_account_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3
`

type ListTransfersForToAccountIdParams struct {
	ToAccountID int32 `json:"to_account_id"`
	Limit       int32 `json:"limit"`
	Offset      int32 `json:"offset"`
}

func (q *Queries) ListTransfersForToAccountId(ctx context.Context, arg ListTransfersForToAccountIdParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, listTransfersForToAccountId, arg.ToAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.PublicID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Amount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}