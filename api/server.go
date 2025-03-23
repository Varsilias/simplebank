package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/varsilias/simplebank/db/sqlc"
	"github.com/varsilias/simplebank/token"
	"github.com/varsilias/simplebank/utils"
)

// Server serves HTTP request to our banking service
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     utils.Config
}

// Response is the wrapper struct sending responses to clients
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// NewServer creates a new HTTP service and sets up routing
func NewServer(store db.Store, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasteoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Authentication
	router.POST("/auth/sign-up", server.registerUser)
	router.POST("/auth/login", server.login)

	// Bank Account Management
	router.GET("/accounts/:public_id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/accounts", server.createAccount)

	// Transfers
	router.POST("/transfers", server.createTransfer)

	server.router = router
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(statusCode int, requestPath string, err error) gin.H {
	return gin.H{
		"timestamp":   time.Now(),
		"path":        requestPath,
		"message":     err.Error(),
		"status_code": statusCode,
	}
}

func successResponse(data any) Response {
	return Response{
		Status:  true,
		Message: "Request successful",
		Data:    data,
	}
}
