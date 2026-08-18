package api

import (
	"context"

	"wallet-api/db/sqlc"
)

// MockStore implements db.Store interface for unit testing HTTP handlers
type MockStore struct {
	db.Store

	CreateUserFn           func(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByEmailFn       func(ctx context.Context, email string) (db.User, error)
	CreateAccountFn        func(ctx context.Context, arg db.CreateAccountParams) (db.Account, error)
	GetAccountFn           func(ctx context.Context, id int64) (db.Account, error)
	GetAccountForUpdateFn func(ctx context.Context, id int64) (db.Account, error)
	ListAccountsFn         func(ctx context.Context, arg db.ListAccountsParams) ([]db.Account, error)
	UpdateAccountBalanceFn func(ctx context.Context, arg db.UpdateAccountBalanceParams) (db.Account, error)
	CreateTransferFn       func(ctx context.Context, arg db.CreateTransferParams) (db.Transfer, error)
	TransferTxFn           func(ctx context.Context, arg db.TransferTxParams) (db.TransferTxResult, error)
}

func (m *MockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, arg)
	}
	return db.User{}, nil
}

func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	if m.GetUserByEmailFn != nil {
		return m.GetUserByEmailFn(ctx, email)
	}
	return db.User{}, nil
}

func (m *MockStore) CreateAccount(ctx context.Context, arg db.CreateAccountParams) (db.Account, error) {
	if m.CreateAccountFn != nil {
		return m.CreateAccountFn(ctx, arg)
	}
	return db.Account{}, nil
}

func (m *MockStore) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	if m.GetAccountFn != nil {
		return m.GetAccountFn(ctx, id)
	}
	return db.Account{}, nil
}

func (m *MockStore) GetAccountForUpdate(ctx context.Context, id int64) (db.Account, error) {
	if m.GetAccountForUpdateFn != nil {
		return m.GetAccountForUpdateFn(ctx, id)
	}
	return db.Account{}, nil
}

func (m *MockStore) ListAccounts(ctx context.Context, arg db.ListAccountsParams) ([]db.Account, error) {
	if m.ListAccountsFn != nil {
		return m.ListAccountsFn(ctx, arg)
	}
	return nil, nil
}

func (m *MockStore) UpdateAccountBalance(ctx context.Context, arg db.UpdateAccountBalanceParams) (db.Account, error) {
	if m.UpdateAccountBalanceFn != nil {
		return m.UpdateAccountBalanceFn(ctx, arg)
	}
	return db.Account{}, nil
}

func (m *MockStore) CreateTransfer(ctx context.Context, arg db.CreateTransferParams) (db.Transfer, error) {
	if m.CreateTransferFn != nil {
		return m.CreateTransferFn(ctx, arg)
	}
	return db.Transfer{}, nil
}

func (m *MockStore) TransferTx(ctx context.Context, arg db.TransferTxParams) (db.TransferTxResult, error) {
	if m.TransferTxFn != nil {
		return m.TransferTxFn(ctx, arg)
	}
	return db.TransferTxResult{}, nil
}
