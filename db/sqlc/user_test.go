package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/varsilias/simplebank/utils"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Firstname, user2.Firstname)
	require.Equal(t, user1.Lastname, user2.Lastname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}

// func TestGetAccountByUserId(t *testing.T) {
// 	account1 := createRandomAccount(t)

// 	account2, err := testQueries.GetAccountByUserId(context.Background(), account1.UserID)
// 	require.NoError(t, err)
// 	require.Equal(t, account1.ID, account2.ID)
// 	require.Equal(t, account1.Currency, account2.Currency)
// 	require.Equal(t, account1.Balance, account2.Balance)
// 	require.Equal(t, account1.UserID, account2.UserID)
// 	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

// }

// func TestUpdateAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)

// 	args := UpdateAccountParams{
// 		ID:      account1.ID,
// 		Balance: utils.RandomAmount(),
// 	}

// 	account2, err := testQueries.UpdateAccount(context.Background(), args)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, account2)

// 	require.Equal(t, account1.ID, account2.ID)
// 	require.Equal(t, account1.Currency, account2.Currency)
// 	require.Equal(t, args.Balance, account2.Balance)
// 	require.Equal(t, account1.UserID, account2.UserID)
// 	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

// }

// func TestDeleteAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	err := testQueries.DeleteAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, account2)
// }

// func TestListAccount(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomAccount(t)
// 	}

// 	args := ListAccountsParams{
// 		Limit:  5,
// 		Offset: 5,
// 	}

// 	accounts, err := testQueries.ListAccounts(context.Background(), args)
// 	require.NoError(t, err)
// 	require.Len(t, accounts, 5)

// 	for _, account := range accounts {
// 		require.NotEmpty(t, account)
// 	}

// }

func createRandomUser(t *testing.T) User {
	// user := createTestUser()
	pass, err := utils.HashPassword(utils.RandomString(10))
	require.NoError(t, err)
	args := CreateUserParams{
		PublicID:  utils.RandomString(26),
		Firstname: utils.RandomString(12),
		Lastname:  utils.RandomString(12),
		Email:     utils.RandomEmail(),
		Password:  pass.HashedPassword,
		Salt:      pass.Salt,
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Firstname, user.Firstname)
	require.Equal(t, args.Lastname, user.Lastname)
	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.Password, user.Password)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}
