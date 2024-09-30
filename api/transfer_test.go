package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/alikhanMuslim/simpleBank/db/mock"
	db "github.com/alikhanMuslim/simpleBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferApi(t *testing.T) {
	amount := 10
	user, _ := RandomUser(t)
	fromAccount := randomAccount(user.Username)
	fromAccount.Currency = "USD"

	user1, _ := RandomUser(t)
	toAccount := randomAccount(user1.Username)
	toAccount.Currency = "USD"

	currency := "USD"

	testcases := []struct {
		name          string
		body          gin.H
		buildStabs    func(store *mockdb.MockStore)
		checkresponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(toAccount, nil)

				arg := db.TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        int64(amount),
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"from_account_id": -92319,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(toAccount, nil)

				arg := db.TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					Amount:        int64(amount),
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(db.TransferTxResults{}, sql.ErrConnDone)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"from_account_id": int64(1001010),
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(int64(1001010))).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Currency is mismatch",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStabs: func(store *mockdb.MockStore) {
				fromAccount.Currency = "KZT"
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError 2",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          amount,
				"currency":        currency,
			},
			buildStabs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkresponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStabs(store)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()
			url := "/transfers"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

		})
	}
}
