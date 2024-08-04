package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Store provide all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

type TransferTxParams struct {
	FromAccountID int32 `json:"from_account_id"`
	ToAccountID   int32 `json:"to_account_id"`
	Amount        int32 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs money transfer from one account to another
// It creates a transfer record, add account entries and update accounts' balance
// Within a single transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)

		if err != nil {
			return err
		}

		result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)

		if err != nil {
			return err
		}

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        int64(arg.Amount),
			PublicID:      uuid.New().String(),
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:   arg.FromAccountID,
			Amount:      -int64(arg.Amount),
			Type:        EntryTypeDEBIT,
			PublicID:    uuid.New().String(),
			LastBalance: result.FromAccount.Balance,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID:   arg.ToAccountID,
			Amount:      int64(arg.Amount),
			Type:        EntryTypeCREDIT,
			PublicID:    uuid.New().String(),
			LastBalance: result.ToAccount.Balance,
		})

		if err != nil {
			return err
		}

		// TODO: Update account balances'
		return nil
	})

	return result, err
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}
