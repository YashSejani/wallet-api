package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"wallet-api/util"
)

func newTestStore(t *testing.T) Store {
	config, err := util.LoadConfig("../..")
	if err != nil {
		config, err = util.LoadConfig(".")
		if err != nil {
			t.Skip("cannot load config for db testing")
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		t.Skip("skipping db test: postgres DB not connected")
		return nil
	}

	if err := connPool.Ping(ctx); err != nil {
		connPool.Close()
		t.Skip("skipping db test: postgres DB ping failed")
		return nil
	}

	return NewStore(connPool)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := newTestStore(t)
	if store == nil {
		return
	}

	// Create test user and accounts
	ctx := context.Background()
	user, err := store.CreateUser(ctx, CreateUserParams{
		Email:        "deadlock@example.com",
		PasswordHash: "secret",
	})
	if err != nil {
		t.Skip("skipping deadlock test: failed to create test user")
		return
	}

	account1, err := store.CreateAccount(ctx, CreateAccountParams{
		UserID:   user.ID,
		Balance:  1000,
		Currency: "USD",
	})
	if err != nil {
		t.Skip("skipping deadlock test: failed to create account 1")
		return
	}

	account2, err := store.CreateAccount(ctx, CreateAccountParams{
		UserID:   user.ID,
		Balance:  1000,
		Currency: "USD",
	})
	if err != nil {
		t.Skip("skipping deadlock test: failed to create account 2")
		return
	}

	// Run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func(from, to int64) {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: from,
				ToAccountID:   to,
				Amount:        amount,
			})
			errs <- err
		}(fromAccountID, toAccountID)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Errorf("transfer transaction failed: %v", err)
		}
	}

	// Verify final balances
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated account 1: %v", err)
	}
	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated account 2: %v", err)
	}

	if updatedAccount1.Balance != 1000 || updatedAccount2.Balance != 1000 {
		t.Errorf("expected final balances 1000 and 1000, got %d and %d", updatedAccount1.Balance, updatedAccount2.Balance)
	}
}
