package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println("<<before", account1.Balance, account2.Balance)
	errs := make(chan error)
	results := make(chan TransferTxResults)

	n := 4
	amount := 10

	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        int64(amount),
	}

	for i := 0; i < n; i++ {
		go func(i int) {
			result, err := store.TransferTx(context.Background(), arg)
			errs <- err
			results <- result
		}(i)
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
		require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
		require.Equal(t, transfer.Amount, arg.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		FromEntry := result.FromEntry
		require.NotEmpty(t, FromEntry)
		require.Equal(t, FromEntry.AccountID, arg.FromAccountID)
		require.Equal(t, FromEntry.Amount, -arg.Amount)
		require.NotZero(t, FromEntry.ID)
		require.NotZero(t, FromEntry.CreatedAt)

		_, err = store.GetEntrie(context.Background(), FromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, ToEntry.AccountID, arg.ToAccountID)
		require.Equal(t, ToEntry.Amount, arg.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = store.GetEntrie(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)
		require.Equal(t, fromAccount.Owner, account1.Owner)
		require.Equal(t, fromAccount.Currency, account1.Currency)
		require.WithinDuration(t, fromAccount.CreatedAt, account1.CreatedAt, time.Second)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.Owner, account2.Owner)
		require.Equal(t, toAccount.Currency, account2.Currency)
		require.WithinDuration(t, toAccount.CreatedAt, account2.CreatedAt, time.Second)

		fmt.Println("<<tx", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%arg.Amount == 0)

		k := int(diff1 / arg.Amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	fmt.Println("<<after", updatedAccount1.Balance, updatedAccount2.Balance)

	bal := n * amount
	require.Equal(t, account1.Balance-int64(bal), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(bal), updatedAccount2.Balance)

}

func TestTransferTXDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println("<<before", account1.Balance, account2.Balance)
	errs := make(chan error)

	n := 10
	amount := int64(10)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func(i int) {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	fmt.Println("<<after", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
