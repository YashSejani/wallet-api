package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wallet-api/db/sqlc"
	"wallet-api/util"
)

func TestCreateTransferAPI(t *testing.T) {
	userID := int64(1)
	otherUserID := int64(2)
	fromAccountID := int64(10)
	toAccountID := int64(20)

	secret := "0123456789abcdef0123456789abcdef"
	config := util.Config{JWTSecret: secret, TokenDuration: 15 * time.Minute}
	token, _ := util.GenerateToken(userID, secret, 15*time.Minute)

	testCases := []struct {
		name          string
		body          any
		setupAuth     func(req *http.Request)
		buildMock     func(mock *MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: transferRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        100,
				Currency:      "USD",
			},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.GetAccountFn = func(ctx context.Context, id int64) (db.Account, error) {
					if id == fromAccountID {
						return db.Account{ID: fromAccountID, UserID: userID, Balance: 500, Currency: "USD"}, nil
					}
					return db.Account{ID: toAccountID, UserID: otherUserID, Balance: 200, Currency: "USD"}, nil
				}
				mock.TransferTxFn = func(ctx context.Context, arg db.TransferTxParams) (db.TransferTxResult, error) {
					return db.TransferTxResult{
						Transfer:    db.Transfer{ID: 1, FromAccountID: arg.FromAccountID, ToAccountID: arg.ToAccountID, Amount: arg.Amount},
						FromAccount: db.Account{ID: fromAccountID, UserID: userID, Balance: 400, Currency: "USD"},
						ToAccount:   db.Account{ID: toAccountID, UserID: otherUserID, Balance: 300, Currency: "USD"},
					}, nil
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
				}
			},
		},
		{
			name: "Unauthorized",
			body: transferRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        100,
				Currency:      "USD",
			},
			setupAuth: func(req *http.Request) {},
			buildMock: func(mock *MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", rec.Code)
				}
			},
		},
		{
			name: "NegativeAmount",
			body: transferRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        -50,
				Currency:      "USD",
			},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("expected status 400, got %d", rec.Code)
				}
			},
		},
		{
			name: "SameAccount",
			body: transferRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   fromAccountID,
				Amount:        100,
				Currency:      "USD",
			},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("expected status 400, got %d", rec.Code)
				}
			},
		},
		{
			name: "SourceForbidden",
			body: transferRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        100,
				Currency:      "USD",
			},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.GetAccountFn = func(ctx context.Context, id int64) (db.Account, error) {
					return db.Account{ID: fromAccountID, UserID: otherUserID, Balance: 500, Currency: "USD"}, nil
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusForbidden {
					t.Errorf("expected status 403, got %d", rec.Code)
				}
			},
		},
		{
			name: "InsufficientBalance",
			body: transferRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        1000,
				Currency:      "USD",
			},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.GetAccountFn = func(ctx context.Context, id int64) (db.Account, error) {
					return db.Account{ID: fromAccountID, UserID: userID, Balance: 500, Currency: "USD"}, nil
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("expected status 400, got %d", rec.Code)
				}
			},
		},
		{
			name: "DestinationNotFound",
			body: transferRequest{
				FromAccountID: fromAccountID,
				ToAccountID:   999,
				Amount:        100,
				Currency:      "USD",
			},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.GetAccountFn = func(ctx context.Context, id int64) (db.Account, error) {
					if id == fromAccountID {
						return db.Account{ID: fromAccountID, UserID: userID, Balance: 500, Currency: "USD"}, nil
					}
					return db.Account{}, errors.New("destination account not found")
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Errorf("expected status 404, got %d", rec.Code)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &MockStore{}
			tc.buildMock(mock)

			server, err := NewServer(config, mock)
			if err != nil {
				t.Fatalf("failed to create server: %v", err)
			}

			data, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			tc.setupAuth(req)
			rec := httptest.NewRecorder()

			server.Router().ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}
