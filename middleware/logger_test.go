package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	})

	loggerHandler := Logger(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	loggerHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("expected body 'ok', got '%s'", rec.Body.String())
	}
}
