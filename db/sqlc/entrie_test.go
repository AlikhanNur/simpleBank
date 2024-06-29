package db

import (
	"context"
	"testing"
	"time"

	"github.com/alikhanMuslim/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntrie(t *testing.T, account Account) Entry {

	arg := CreateEntrieParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entrie, err := testqueries.CreateEntrie(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entrie)

	require.Equal(t, account.ID, entrie.AccountID)
	require.Equal(t, arg.Amount, entrie.Amount)
	require.NotZero(t, entrie.ID)
	require.NotZero(t, entrie.CreatedAt)

	return entrie

}

func TestCreateEntrie(t *testing.T) {
	account := CreateRandomAccount(t)

	CreateRandomEntrie(t, account)
}

func TestGetEntrie(t *testing.T) {
	account := CreateRandomAccount(t)

	entrie1 := CreateRandomEntrie(t, account)

	entrie2, err := testqueries.GetEntrie(context.Background(), entrie1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entrie2)

	require.Equal(t, entrie1.ID, entrie2.ID)
	require.Equal(t, entrie1.AccountID, entrie2.AccountID)
	require.Equal(t, entrie1.Amount, entrie2.Amount)
	require.WithinDuration(t, entrie1.CreatedAt, entrie2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomEntrie(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testqueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)

	for _, entrie := range entries {
		require.NotEmpty(t, entrie)
		require.Equal(t, entrie.AccountID, account.ID)
	}

}
