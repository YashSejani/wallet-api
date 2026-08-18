package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wallet-api/util"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef"
	userID := int64(42)

	validToken, err := util.GenerateToken(userID, secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate valid token: %v", err)
	}

	expiredToken, err := util.GenerateToken(userID, secret, -time.Minute)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	testCases := []struct {
		name          string
		setupAuth     func(req *http.Request)
		checkResponse func(t *testing.T, rec *httptest.ResponseRecorder, contextUserID int64)
	}{
		{
			name: "OK_ValidToken",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder, contextUserID int64) {
				if rec.Code != http.StatusOK {
					t.Errorf("expected status 200, got %d", rec.Code)
				}
				if contextUserID != userID {
					t.Errorf("expected context userID %d, got %d", userID, contextUserID)
				}
			},
		},
		{
			name: "MissingHeader",
			setupAuth: func(req *http.Request) {
				// No Authorization header set
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder, contextUserID int64) {
				if rec.Code != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", rec.Code)
				}
			},
		},
		{
			name: "InvalidHeaderFormat",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Basic invalidtoken")
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder, contextUserID int64) {
				if rec.Code != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", rec.Code)
				}
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", expiredToken))
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder, contextUserID int64) {
				if rec.Code != http.StatusUnauthorized {
					t.Errorf("expected status 401, got %d", rec.Code)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var capturedUserID int64
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if val, ok := r.Context().Value(UserIDKey).(int64); ok {
					capturedUserID = val
				}
				w.WriteHeader(http.StatusOK)
			})

			authMiddleware := Auth(secret)
			handlerToTest := authMiddleware(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			tc.setupAuth(req)
			rec := httptest.NewRecorder()

			handlerToTest.ServeHTTP(rec, req)
			tc.checkResponse(t, rec, capturedUserID)
		})
	}
}
