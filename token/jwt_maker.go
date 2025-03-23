package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil

}

// CreateToken creates a new token for a specific user
func (maker *JWTMaker) CreateToken(publicID string, duration time.Duration) (string, error) {
	payload, err := NewPayload(publicID, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.secretKey))
}

// VerifyToken verifies the authenticity of a token
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	fmt.Printf("Error Type: %T", err)
	if err != nil {

		switch {
		case errors.Is(err, jwt.ErrTokenInvalidClaims):
			return nil, ErrInvalidToken
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, ErrInvalidToken
		case errors.Is(err, jwt.ErrTokenUnverifiable):
			return nil, ErrInvalidToken
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrInvalidToken
		default:
			return nil, ErrExpiredToken
		}
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
