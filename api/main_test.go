package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"wallet-api/db/sqlc"
	"wallet-api/util"
)

type mockStore struct {
	db.Store
	createUserFn    func(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	getUserByEmailFn func(ctx context.Context, email string) (db.User, error)
}

func (m *mockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, arg)
	}
	return db.User{}, nil
}

func (m *mockStore) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(ctx, email)
	}
	return db.User{}, nil
}

func TestCreateUserAPI(t *testing.T) {
	config := util.Config{
		JWTSecret:     "test_secret_12345678901234567890123456789012",
		TokenDuration: 15 * time.Minute,
	}

	mock := &mockStore{
		createUserFn: func(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
			return db.User{
				ID:           1,
				Email:        arg.Email,
				PasswordHash: arg.PasswordHash,
				CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			}, nil
		},
	}

	server, err := NewServer(config, mock)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	body, _ := json.Marshal(createUserRequest{
		Email:    "test@example.com",
		Password: "secretpassword",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
}
