package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(12)

	hashData, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashData)

	isMatch, err := VerifyPassword(password, hashData.HashedPassword, hashData.Salt)
	require.NoError(t, err)
	require.True(t, isMatch)

	wrongPassword := RandomString(10)

	isMatch, err = VerifyPassword(wrongPassword, hashData.HashedPassword, hashData.Salt)
	require.NoError(t, err)
	require.False(t, isMatch)

}
