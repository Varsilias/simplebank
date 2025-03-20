package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/varsilias/simplebank/db/sqlc"
)

// Server serves HTTP request to our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Response is the wrapper struct sending responses to clients
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// NewServer creates a new HTTP service and sets up routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// Authentication
	router.POST("/auth/sign-up", server.registerUser)

	// Bank Account Management
	router.GET("/accounts/:public_id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	// Transfers
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(statusCode int, requestPath string, err error) gin.H {
	return gin.H{
		"timestamp":  time.Now(),
		"path":       requestPath,
		"message":    err.Error(),
		"statusCode": statusCode,
	}
}

func successResponse(data any) Response {
	return Response{
		Status:  true,
		Message: "Request successful",
		Data:    data,
	}
}
