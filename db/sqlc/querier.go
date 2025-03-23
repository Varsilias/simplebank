// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccount(ctx context.Context, id int32) error
	DeleteEntry(ctx context.Context, id int32) error
	DeleteTransfer(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, id int32) error
	GetAccount(ctx context.Context, id int32) (Account, error)
	GetAccountByPublicId(ctx context.Context, publicID string) (Account, error)
	GetAccountByUserId(ctx context.Context, userID int32) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int32) (Account, error)
	GetAllUserAccounts(ctx context.Context, userID int32) ([]Account, error)
	GetEntry(ctx context.Context, id int32) (Entry, error)
	GetEntryByAccountId(ctx context.Context, accountID int32) (Entry, error)
	GetTransfer(ctx context.Context, id int32) (Transfer, error)
	GetTransferByFromAccountId(ctx context.Context, fromAccountID int32) (Transfer, error)
	GetTransferByToAccountId(ctx context.Context, toAccountID int32) (Transfer, error)
	GetUser(ctx context.Context, id int32) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByPublicID(ctx context.Context, publicID string) (User, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListEntriesForAccountId(ctx context.Context, arg ListEntriesForAccountIdParams) ([]Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	ListTransfersForFromAccountId(ctx context.Context, arg ListTransfersForFromAccountIdParams) ([]Transfer, error)
	ListTransfersForToAccountId(ctx context.Context, arg ListTransfersForToAccountIdParams) ([]Transfer, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
