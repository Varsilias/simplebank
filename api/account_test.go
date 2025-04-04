package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mockdb "github.com/varsilias/simplebank/db/mock"
	db "github.com/varsilias/simplebank/db/sqlc"
	"github.com/varsilias/simplebank/token"
	"github.com/varsilias/simplebank/utils"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := createRandomUser(t)
	account := createRandomAccount(user.ID)

	testCases := []struct {
		name            string
		accountPublicID string
		setupAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs      func(store *mockdb.MockStore)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:            "Ok",
			accountPublicID: account.PublicID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByPublicId(gomock.Any(), gomock.Eq(account.PublicID)).Times(1).Return(account, nil)
				store.EXPECT().GetUserByPublicID(gomock.Any(), gomock.Eq(user.PublicID)).Times(1).Return(user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorisation(t, request, tokenMaker, authorisationTypeBearer, user.PublicID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:            "NotFound",
			accountPublicID: account.PublicID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByPublicID(gomock.Any(), gomock.Eq(user.PublicID)).Times(1).Return(user, nil)
				store.EXPECT().GetAccountByPublicId(gomock.Any(), gomock.Eq(account.PublicID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorisation(t, request, tokenMaker, authorisationTypeBearer, user.PublicID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:            "InternalServerError",
			accountPublicID: account.PublicID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByPublicID(gomock.Any(), gomock.Eq(user.PublicID)).Times(1).Return(user, nil)
				store.EXPECT().GetAccountByPublicId(gomock.Any(), gomock.Eq(account.PublicID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorisation(t, request, tokenMaker, authorisationTypeBearer, user.PublicID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:            "BadRequest",
			accountPublicID: "not-a-valid-uuid",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByPublicId(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorisation(t, request, tokenMaker, authorisationTypeBearer, user.PublicID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%s", tc.accountPublicID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}

}

func createRandomAccount(userID int32) db.Account {
	return db.Account{
		ID:       int32(utils.RandomInt(1, 1000)),
		PublicID: uuid.New().String(),
		UserID:   userID,
		Balance:  utils.RandomAmount(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	fmt.Println("Raw JSON:", string(data))

	var response Response
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	dataMap, ok := response.Data.(map[string]any)
	require.True(t, ok, "data is not a map")

	gotAccount := db.Account{
		ID:       int32(dataMap["id"].(float64)),
		PublicID: dataMap["public_id"].(string),
		UserID:   int32(dataMap["user_id"].(float64)),
		Balance:  int64(dataMap["balance"].(float64)),
		Currency: dataMap["currency"].(string),
	}

	require.Equal(t, account.ID, gotAccount.ID)
	require.Equal(t, account.PublicID, gotAccount.PublicID)
	require.Equal(t, account.UserID, gotAccount.UserID)
	require.Equal(t, account.Balance, gotAccount.Balance)
	require.Equal(t, account.Currency, gotAccount.Currency)
	require.Equal(t, account.IsBlocked, gotAccount.IsBlocked)
}
