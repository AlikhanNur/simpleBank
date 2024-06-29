package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alikhanMuslim/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testqueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	account2, err := testqueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testqueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	err := testqueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testqueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)

	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i <= 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testqueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
