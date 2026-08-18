package api

import (
	"testing"
	"time"

	"wallet-api/util"
)

func TestNewServer(t *testing.T) {
	config := util.Config{
		JWTSecret:     "0123456789abcdef0123456789abcdef",
		TokenDuration: 15 * time.Minute,
	}

	mock := &MockStore{}
	server, err := NewServer(config, mock)

	if err != nil {
		t.Fatalf("expected no error when creating server, got: %v", err)
	}
	if server == nil {
		t.Fatal("expected server to be non-nil")
	}
	if server.Router() == nil {
		t.Fatal("expected server router to be non-nil")
	}
}
