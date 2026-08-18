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

	"github.com/jackc/pgx/v5/pgtype"

	"wallet-api/db/sqlc"
	"wallet-api/util"
)

func TestCreateAccountAPI(t *testing.T) {
	userID := int64(1)
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
			body: createAccountRequest{Currency: "USD"},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.CreateAccountFn = func(ctx context.Context, arg db.CreateAccountParams) (db.Account, error) {
					return db.Account{
						ID:        10,
						UserID:    arg.UserID,
						Balance:   0,
						Currency:  arg.Currency,
						CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
					}, nil
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusCreated {
					t.Errorf("expected status 201, got %d. Body: %s", rec.Code, rec.Body.String())
				}
			},
		},
		{
			name: "Unauthorized",
			body: createAccountRequest{Currency: "USD"},
			setupAuth: func(req *http.Request) {
				// No token set
			},
			buildMock: func(mock *MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", rec.Code)
				}
			},
		},
		{
			name: "InvalidCurrency",
			body: createAccountRequest{Currency: "XYZ"},
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
			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			tc.setupAuth(req)
			rec := httptest.NewRecorder()

			server.Router().ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func TestGetAccountAPI(t *testing.T) {
	userID := int64(1)
	otherUserID := int64(2)
	accountID := int64(10)
	secret := "0123456789abcdef0123456789abcdef"
	config := util.Config{JWTSecret: secret, TokenDuration: 15 * time.Minute}

	token, _ := util.GenerateToken(userID, secret, 15*time.Minute)

	testCases := []struct {
		name          string
		accountID     string
		setupAuth     func(req *http.Request)
		buildMock     func(mock *MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: fmt.Sprintf("%d", accountID),
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.GetAccountFn = func(ctx context.Context, id int64) (db.Account, error) {
					return db.Account{
						ID:        accountID,
						UserID:    userID,
						Balance:   1000,
						Currency:  "USD",
						CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
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
			name:      "AccountNotFound",
			accountID: "999",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.GetAccountFn = func(ctx context.Context, id int64) (db.Account, error) {
					return db.Account{}, errors.New("account not found")
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusNotFound {
					t.Errorf("expected status 404, got %d", rec.Code)
				}
			},
		},
		{
			name:      "ForbiddenUser",
			accountID: fmt.Sprintf("%d", accountID),
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			buildMock: func(mock *MockStore) {
				mock.GetAccountFn = func(ctx context.Context, id int64) (db.Account, error) {
					return db.Account{
						ID:       accountID,
						UserID:   otherUserID, // belongs to user 2
						Balance:  1000,
						Currency: "USD",
					}, nil
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusForbidden {
					t.Errorf("expected status 403, got %d", rec.Code)
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", tc.accountID), nil)
			tc.setupAuth(req)
			rec := httptest.NewRecorder()

			server.Router().ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	userID := int64(1)
	secret := "0123456789abcdef0123456789abcdef"
	config := util.Config{JWTSecret: secret, TokenDuration: 15 * time.Minute}

	token, _ := util.GenerateToken(userID, secret, 15*time.Minute)

	mock := &MockStore{
		ListAccountsFn: func(ctx context.Context, arg db.ListAccountsParams) ([]db.Account, error) {
			return []db.Account{
				{ID: 1, UserID: userID, Balance: 100, Currency: "USD"},
				{ID: 2, UserID: userID, Balance: 200, Currency: "EUR"},
			}, nil
		},
	}

	server, err := NewServer(config, mock)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/accounts?page_id=1&page_size=10", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}
