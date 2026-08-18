package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"wallet-api/db/sqlc"
	"wallet-api/util"
)

func TestCreateUserAPI(t *testing.T) {
	password := "secret123"

	config := util.Config{
		JWTSecret:     "0123456789abcdef0123456789abcdef",
		TokenDuration: 15 * time.Minute,
	}

	testCases := []struct {
		name          string
		body          any
		buildMock     func(mock *MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createUserRequest{
				Email:    "test@example.com",
				Password: password,
			},
			buildMock: func(mock *MockStore) {
				mock.CreateUserFn = func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
					return db.User{
						ID:           1,
						Email:        arg.Email,
						PasswordHash: arg.PasswordHash,
						CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
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
			name: "DuplicateEmail",
			body: createUserRequest{
				Email:    "existing@example.com",
				Password: password,
			},
			buildMock: func(mock *MockStore) {
				mock.CreateUserFn = func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
					return db.User{}, errors.New("duplicate key value violates unique constraint")
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusForbidden {
					t.Errorf("expected status 403, got %d", rec.Code)
				}
			},
		},
		{
			name: "ShortPassword",
			body: createUserRequest{
				Email:    "test@example.com",
				Password: "short",
			},
			buildMock: func(mock *MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("expected status 400, got %d", rec.Code)
				}
			},
		},
		{
			name: "MissingEmail",
			body: createUserRequest{
				Email:    "",
				Password: password,
			},
			buildMock: func(mock *MockStore) {},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("expected status 400, got %d", rec.Code)
				}
			},
		},
		{
			name: "InternalError",
			body: createUserRequest{
				Email:    "test@example.com",
				Password: password,
			},
			buildMock: func(mock *MockStore) {
				mock.CreateUserFn = func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
					return db.User{}, errors.New("database connection failed")
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusInternalServerError {
					t.Errorf("expected status 500, got %d", rec.Code)
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
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.Router().ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	password := "secret123"
	hashedPassword, _ := util.HashPassword(password)

	config := util.Config{
		JWTSecret:     "0123456789abcdef0123456789abcdef",
		TokenDuration: 15 * time.Minute,
	}

	user := db.User{
		ID:           1,
		Email:        "user@example.com",
		PasswordHash: hashedPassword,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	testCases := []struct {
		name          string
		body          any
		buildMock     func(mock *MockStore)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: loginUserRequest{
				Email:    user.Email,
				Password: password,
			},
			buildMock: func(mock *MockStore) {
				mock.GetUserByEmailFn = func(ctx context.Context, email string) (db.User, error) {
					return user, nil
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
				}
			},
		},
		{
			name: "UserNotFound",
			body: loginUserRequest{
				Email:    "nonexistent@example.com",
				Password: password,
			},
			buildMock: func(mock *MockStore) {
				mock.GetUserByEmailFn = func(ctx context.Context, email string) (db.User, error) {
					return db.User{}, errors.New("user not found")
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", rec.Code)
				}
			},
		},
		{
			name: "WrongPassword",
			body: loginUserRequest{
				Email:    user.Email,
				Password: "wrongpassword",
			},
			buildMock: func(mock *MockStore) {
				mock.GetUserByEmailFn = func(ctx context.Context, email string) (db.User, error) {
					return user, nil
				}
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", rec.Code)
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
			req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.Router().ServeHTTP(rec, req)
			tc.checkResponse(t, rec)
		})
	}
}
