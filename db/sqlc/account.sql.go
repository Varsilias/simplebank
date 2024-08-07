// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: account.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  public_id, user_id, balance, currency
) VALUES (
  $1, $2, $3, $4
) RETURNING id, public_id, is_blocked, blocked_at, created_at, updated_at, deleted_at, user_id, balance, currency
`

type CreateAccountParams struct {
	PublicID string `json:"public_id"`
	UserID   int32  `json:"user_id"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount,
		arg.PublicID,
		arg.UserID,
		arg.Balance,
		arg.Currency,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.IsBlocked,
		&i.BlockedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.Balance,
		&i.Currency,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, public_id, is_blocked, blocked_at, created_at, updated_at, deleted_at, user_id, balance, currency FROM accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int32) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.IsBlocked,
		&i.BlockedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.Balance,
		&i.Currency,
	)
	return i, err
}

const getAccountByUserId = `-- name: GetAccountByUserId :one
SELECT id, public_id, is_blocked, blocked_at, created_at, updated_at, deleted_at, user_id, balance, currency FROM accounts
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetAccountByUserId(ctx context.Context, userID int32) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountByUserId, userID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.IsBlocked,
		&i.BlockedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.Balance,
		&i.Currency,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, public_id, is_blocked, blocked_at, created_at, updated_at, deleted_at, user_id, balance, currency FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.IsBlocked,
			&i.BlockedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.UserID,
			&i.Balance,
			&i.Currency,
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

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING id, public_id, is_blocked, blocked_at, created_at, updated_at, deleted_at, user_id, balance, currency
`

type UpdateAccountParams struct {
	ID      int32  `json:"id"`
	Balance string `json:"balance"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.IsBlocked,
		&i.BlockedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.Balance,
		&i.Currency,
	)
	return i, err
}
