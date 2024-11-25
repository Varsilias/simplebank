package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/varsilias/simplebank/db/sqlc"
	"time"
)

// Server serves HTTP request to our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Response is the wrapper struct sending responses to clients
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// NewServer creates a new HTTP service and sets up routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Authentication
	router.POST("/auth/sign-up", server.registerUser)

	// Bank Account Management
	router.GET("/accounts/:public_id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

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
