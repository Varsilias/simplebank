package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/varsilias/simplebank/db/sqlc"
	"github.com/varsilias/simplebank/utils"
	"log"
	"net/http"
	"time"
)

type createUserRequest struct {
	Firstname string `json:"firstname" binding:"required,min=3"`
	Lastname  string `json:"lastname" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Currency  string `json:"currency" binding:"required,oneof=USD EUR GBP NGN"`
}

type createUserResponse struct {
	ID                       int32                  `json:"id"`
	PublicID                 string                 `json:"public_id"`
	IsBlocked                bool                   `json:"is_blocked"`
	BlockedAt                *string                `json:"blocked_at"`
	CreatedAt                time.Time              `json:"created_at"`
	UpdatedAt                time.Time              `json:"updated_at"`
	DeletedAt                *string                `json:"deleted_at"`
	Firstname                string                 `json:"firstname"`
	Lastname                 string                 `json:"lastname"`
	Email                    string                 `json:"email"`
	Password                 string                 `json:"-"`
	Salt                     string                 `json:"-"`
	SecurityToken            *string                `json:"security_token"`
	EmailConfirmed           bool                   `json:"email_confirmed"`
	SecurityTokenRequestedAt *string                `json:"security_token_requested_at"`
	AccountDetail            *createAccountResponse `json:"account_detail"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(http.StatusBadRequest, ctx.Request.URL.Path, err))
		return
	}

	userExists, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("Error Getting user: ", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, err))
			return
		}
	}

	var emptyUser db.User
	if userExists != emptyUser {
		ctx.JSON(http.StatusConflict, errorResponse(http.StatusConflict, ctx.Request.URL.Path, errors.New("user already exists")))
		return
	}

	if !utils.IsValidPassword(req.Password) {
		ctx.JSON(http.StatusBadRequest, errorResponse(http.StatusBadRequest, ctx.Request.URL.Path, errors.New("password must be at least 8 characters long, contain 1 uppercase, 1 lowercase, 1 special character and a numbe")))
		return
	}

	passwordData, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, errors.New("something went wrong, we are fixing it")))
		return
	}

	arg := db.CreateUserParams{
		PublicID:  uuid.New().String(),
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Password:  passwordData.HashedPassword,
		Salt:      passwordData.Salt,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, err))
		return
	}

	account, err := server.createAccountWithArgs(ctx, createAccountRequest{
		UserID:   user.ID,
		Currency: req.Currency,
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(http.StatusInternalServerError, ctx.Request.URL.Path, errors.New("something went wrong, we are fixing it")))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(toUserResponse(user, account)))
}

// Convert User model to the Response struct
func toUserResponse(user db.User, account *createAccountResponse) *createUserResponse {
	var blockedAt, deletedAt, securityTokenRequestedAt *string

	// Convert nullable time fields
	if user.BlockedAt.Valid {
		t := user.BlockedAt.Time.Format(time.RFC3339)
		blockedAt = &t
	}
	if user.DeletedAt.Valid {
		t := user.DeletedAt.Time.Format(time.RFC3339)
		deletedAt = &t
	}
	if user.SecurityTokenRequestedAt.Valid {
		t := user.SecurityTokenRequestedAt.Time.Format(time.RFC3339)
		securityTokenRequestedAt = &t
	}

	// Convert security token
	var securityToken *string
	if user.SecurityToken.Valid {
		securityToken = &user.SecurityToken.String
	}

	return &createUserResponse{
		ID:                       user.ID,
		PublicID:                 user.PublicID,
		IsBlocked:                user.IsBlocked.Bool,
		BlockedAt:                blockedAt,
		CreatedAt:                user.CreatedAt,
		UpdatedAt:                user.UpdatedAt,
		DeletedAt:                deletedAt,
		Firstname:                user.Firstname,
		Lastname:                 user.Lastname,
		Email:                    user.Email,
		Password:                 user.Password,
		Salt:                     user.Salt,
		SecurityToken:            securityToken,
		EmailConfirmed:           user.EmailConfirmed.Bool,
		SecurityTokenRequestedAt: securityTokenRequestedAt,
		AccountDetail:            account,
	}
}