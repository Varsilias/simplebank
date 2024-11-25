package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/varsilias/simplebank/db/sqlc"
	"log"
	"net/http"
	"time"
)

type createAccountResponse struct {
	ID        int32     `json:"id"`
	PublicID  string    `json:"public_id"`
	IsBlocked bool      `json:"is_blocked"`
	BlockedAt *string   `json:"blocked_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *string   `json:"deleted_at"`
	UserID    int32     `json:"user_id"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
}

// createAccountRequest is the type for creating a new account
type createAccountRequest struct {
	UserID   int32  `json:"user_id"`
	Currency string `json:"currency"`
}

func (server *Server) createAccountWithArgs(ctx *gin.Context, createAccountArgs createAccountRequest) (*createAccountResponse, error) {
	var account db.Account
	accountExists, err := server.store.GetAccountByUserId(ctx, createAccountArgs.UserID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("Error getting account: ", err)
			return nil, err
		}
	}

	if accountExists != account {
		return toAccountResponse(accountExists), nil
	}

	args := db.CreateAccountParams{
		PublicID: uuid.New().String(),
		UserID:   createAccountArgs.UserID,
		Balance:  0,
		Currency: createAccountArgs.Currency,
	}

	account, err = server.store.CreateAccount(ctx, args)

	return toAccountResponse(account), err
}

type getAccountRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(http.StatusBadRequest, ctx.Request.URL.Path, err))
	}

	account, err := server.store.GetAccountByPublicId(ctx, req.PublicID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(http.StatusNotFound, ctx.Request.URL.Path, errors.New("account not found")))
			return
		}
		log.Println("Error getting account: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(toAccountResponse(account)))
}

type listAccountsRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"pageSize" binding:"required,min=10,max=20"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(http.StatusBadRequest, ctx.Request.URL.Path, err))
		return
	}

	args := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, err))
	}
	accountList := make([]*createAccountResponse, 0)
	for _, account := range accounts {
		accountList = append(accountList, toAccountResponse(account))
	}
	ctx.JSON(http.StatusOK, successResponse(accountList))
}

func toAccountResponse(account db.Account) *createAccountResponse {
	var blockedAt, deletedAt *string

	if account.BlockedAt.Valid {
		t := account.BlockedAt.Time.Format(time.RFC3339)
		blockedAt = &t
	}
	if account.DeletedAt.Valid {
		t := account.DeletedAt.Time.Format(time.RFC3339)
		deletedAt = &t
	}

	return &createAccountResponse{
		ID:        account.ID,
		PublicID:  account.PublicID,
		IsBlocked: account.IsBlocked.Bool,
		BlockedAt: blockedAt,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		DeletedAt: deletedAt,
		UserID:    account.UserID,
		Balance:   account.Balance,
		Currency:  account.Currency,
	}
}
