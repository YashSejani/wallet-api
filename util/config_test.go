package util

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("..")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if config.DBSource == "" {
		t.Errorf("expected DBSource to be non-empty")
	}

	if config.ServerAddress == "" {
		t.Errorf("expected ServerAddress to be non-empty")
	}

	if config.JWTSecret == "" {
		t.Errorf("expected JWTSecret to be non-empty")
	}

	if config.TokenDuration != 15*time.Minute {
		t.Errorf("expected TokenDuration 15m, got %v", config.TokenDuration)
	}
}
