package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/varsilias/simplebank/db/sqlc"
	"github.com/varsilias/simplebank/token"
)

type createTransferTxResponse struct {
	Transfer    createTransferResponse `json:"transfer"`
	FromAccount createAccountResponse  `json:"from_account"`
	ToAccount   createAccountResponse  `json:"to_account"`
	FromEntry   createEntryResponse    `json:"from_entry"`
	ToEntry     createEntryResponse    `json:"to_entry"`
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(http.StatusBadRequest, ctx.Request.URL.Path, err))
		return
	}

	authPayload := ctx.MustGet(authorisationKey).(*token.Payload)
	user, err := server.store.GetUserByPublicID(ctx, authPayload.PublicID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("user with id [%s] not found", authPayload.PublicID)
			ctx.JSON(http.StatusNotFound, errorResponse(http.StatusNotFound, ctx.Request.URL.Path, err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, int32(req.FromAccountID), req.Currency)
	if !valid {
		return
	}

	if fromAccount.UserID != user.ID {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(http.StatusForbidden, ctx.Request.URL.Path, err))
		return
	}

	_, valid = server.validAccount(ctx, int32(req.ToAccountID), req.Currency)
	if !valid {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: int32(req.FromAccountID),
		ToAccountID:   int32(req.ToAccountID),
		Amount:        int32(req.Amount),
	}

	result, err := server.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(toTransferTxResponse(result)))
}

func (server *Server) validAccount(ctx *gin.Context, accountID int32, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(http.StatusNotFound, ctx.Request.URL.Path, err))
			return account, false
		}

		log.Println("Error getting account: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(http.StatusBadRequest, ctx.Request.URL.Path, err))
		return account, false
	}

	return account, true
}

func toTransferTxResponse(transferTxResult db.TransferTxResult) *createTransferTxResponse {
	fromAccountResponse := toAccountResponse(transferTxResult.FromAccount)
	toAccountResponse := toAccountResponse(transferTxResult.ToAccount)
	fromEntryResponse := toEntryResponse(transferTxResult.FromEntry)
	toEntryResponse := toEntryResponse(transferTxResult.ToEntry)
	transferResponse := toTransferResponse(transferTxResult.Transfer)

	return &createTransferTxResponse{
		Transfer:    *transferResponse,
		FromAccount: *fromAccountResponse,
		ToAccount:   *toAccountResponse,
		FromEntry:   *fromEntryResponse,
		ToEntry:     *toEntryResponse,
	}
}

type createTransferResponse struct {
	ID            int32     `json:"id"`
	PublicID      string    `json:"public_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *string   `json:"deleted_at"`
	FromAccountID int32     `json:"from_acount_id"`
	ToAccountID   int32     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
}

func toTransferResponse(transfer db.Transfer) *createTransferResponse {
	var deletedAt *string

	if transfer.DeletedAt.Valid {
		t := transfer.DeletedAt.Time.Format(time.RFC3339)
		deletedAt = &t
	}

	return &createTransferResponse{
		ID:            transfer.ID,
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		PublicID:      transfer.PublicID,
		CreatedAt:     transfer.CreatedAt,
		UpdatedAt:     transfer.UpdatedAt,
		DeletedAt:     deletedAt,
		Amount:        transfer.Amount,
	}
}

type createEntryResponse struct {
	ID          int32     `json:"id"`
	PublicID    string    `json:"public_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *string   `json:"deleted_at"`
	AccountID   int32     `json:"account_id"`
	Amount      int64     `json:"amount"`
	LastBalance int64     `json:"last_balance"`
	Type        string    `json:"type"`
}

func toEntryResponse(entry db.Entry) *createEntryResponse {
	var deletedAt *string

	if entry.DeletedAt.Valid {
		t := entry.DeletedAt.Time.Format(time.RFC3339)
		deletedAt = &t
	}

	return &createEntryResponse{
		ID:          entry.ID,
		AccountID:   entry.AccountID,
		PublicID:    entry.PublicID,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
		DeletedAt:   deletedAt,
		Amount:      entry.Amount,
		Type:        string(entry.Type),
		LastBalance: entry.LastBalance,
	}
}
