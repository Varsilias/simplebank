package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mockdb "github.com/varsilias/simplebank/db/mock"
	db "github.com/varsilias/simplebank/db/sqlc"
	"github.com/varsilias/simplebank/utils"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	args     db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	args, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	isMatch, err := utils.VerifyPassword(e.password, args.Password, args.Salt)

	if err != nil {
		return false
	}

	e.args.Password = args.Password

	// Check everything except PublicID
	return args.Firstname == e.args.Firstname &&
		args.Lastname == e.args.Lastname &&
		args.Email == e.args.Email &&
		isMatch
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches args %v and password %v", e.args, e.password)
}

func EqCreateUserParams(args db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{args, password}
}

func TestRegisterUserAPI(t *testing.T) {
	user, password := createRandomUser(t)
	currency := utils.RandomCurrency()
	account := createRandomAccount()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     user.Email,
				"password":  password,
				"currency":  currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Firstname: user.Firstname,
					Lastname:  user.Lastname,
					Email:     user.Email,
					Salt:      user.Salt,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(args, password)).Times(1).Return(user, nil)
				store.EXPECT().GetAccountByUserId(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, nil)
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(_ any, arg db.CreateAccountParams) (db.Account, error) {
						require.Equal(t, user.ID, arg.UserID)
						require.Equal(t, currency, arg.Currency)
						require.Equal(t, int64(0), arg.Balance)
						require.NotEmpty(t, arg.PublicID)

						return account, nil
					})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				log.Printf("Recorder Body: %v", recorder.Body)
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user, account)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/auth/sign-up"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomPassword(10) // should be minimum of 8
	fmt.Printf("Generated Password: %s\n", password)
	pass, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:        int32(utils.RandomInt(1, 1000)),
		PublicID:  uuid.New().String(),
		Firstname: utils.RandomString(10),
		Lastname:  utils.RandomString(10),
		Email:     utils.RandomEmail(),
		Password:  pass.HashedPassword,
		Salt:      pass.Salt,
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	fmt.Println("Raw JSON:", string(data))

	var response Response
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	dataMap, ok := response.Data.(map[string]any)
	require.True(t, ok, "data is not a map")

	gotUser := db.User{
		ID:        int32(dataMap["id"].(float64)),
		Firstname: dataMap["firstname"].(string),
		Lastname:  dataMap["lastname"].(string),
		Email:     dataMap["email"].(string),
	}

	require.Equal(t, user.ID, gotUser.ID)
	// require.Equal(t, user.PublicID, gotUser.PublicID)
	require.Equal(t, user.Firstname, gotUser.Firstname)
	require.Equal(t, user.Lastname, gotUser.Lastname)
	require.Equal(t, user.Email, gotUser.Email)

	// Check account details
	accountDetails, ok := dataMap["account_detail"].(map[string]any)
	require.True(t, ok, "account_detail is not a map")

	if accountDetails != nil {
		gotAccount := db.Account{
			ID:       int32(accountDetails["id"].(float64)),
			PublicID: accountDetails["public_id"].(string),
			UserID:   int32(accountDetails["user_id"].(float64)),
			Balance:  int64(accountDetails["balance"].(float64)),
			Currency: accountDetails["currency"].(string),
		}

		require.Equal(t, account.ID, gotAccount.ID)
		require.Equal(t, account.PublicID, gotAccount.PublicID)
		require.Equal(t, account.UserID, gotAccount.UserID)
		require.Equal(t, account.Balance, gotAccount.Balance)
		require.Equal(t, account.Currency, gotAccount.Currency)
	}
}
