package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/utils"
	"testing"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)

	// 创建账号
	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()

	account1.Currency = utils.USD
	account2.Currency = utils.USD
	account3.Currency = utils.EUR

	// 正常流程
	/*
		1. 验证from 用户
		2. 验证to 用户
		3. 创建transfer记录
		4. 创建的两天余额记录
	*/
	var testcases = []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, pgx.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, pgx.ErrNoRows)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)

				arg := db.TransferTxParams{
					FromAccountID: account3.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account3.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "XYZ",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          -amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, sql.ErrConnDone)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建ctril
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 生成gomock store
			store := mockdb.NewMockStore(ctrl)

			// 开启mockdb预构建
			tc.buildStubs(store)

			// 生成 gin server
			server := NewServer(store)

			url := fmt.Sprintf("/transfers")

			jsonData, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			// 创建响应接收器
			recorder := httptest.NewRecorder()
			// 获取响应
			server.router.ServeHTTP(recorder, request)

			// 检查响应
			tc.checkResponse(t, recorder)

		})
	}

}

func randomTransfer(fromAccount, toAccount db.Account) db.Transfer {
	return db.Transfer{
		ID:            utils.RandomInt(1, 1000),
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        utils.RandomMoney(),
	}
}
