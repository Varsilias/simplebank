package token

import "time"

// Maker is the interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific user
	CreateToken(publicID string, duration time.Duration) (string, error)

	// VerifyToken verifies the authenticity of a token
	VerifyToken(token string) (*Payload, error)
}
