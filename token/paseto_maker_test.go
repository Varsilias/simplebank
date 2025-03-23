package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/varsilias/simplebank/utils"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasteoMaker(utils.RandomPassword(16))
	require.NoError(t, err)

	publicID := uuid.New().String()
	duration := time.Minute
	issueAt := time.Now()
	expiredAt := issueAt.Add(duration)

	token, err := maker.CreateToken(publicID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, publicID, payload.PublicID)
	require.Equal(t, publicID, payload.PublicID)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasteoMaker(utils.RandomPassword(16))
	require.NoError(t, err)

	publicID := uuid.New().String()
	duration := time.Minute

	token, err := maker.CreateToken(publicID, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
