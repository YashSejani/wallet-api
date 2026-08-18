package util

import (
	"testing"
	"time"
)

func TestJWTToken(t *testing.T) {
	userID := int64(123)
	secret := "0123456789abcdef0123456789abcdef"
	duration := 15 * time.Minute

	// Test 1: Valid Token Generation & Validation
	token, err := GenerateToken(userID, secret, duration)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if len(token) == 0 {
		t.Fatal("expected non-empty token string")
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("expected userID %d, got %d", userID, claims.UserID)
	}

	// Test 2: Expired Token Validation
	expiredToken, err := GenerateToken(userID, secret, -time.Minute)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}
	_, err = ValidateToken(expiredToken, secret)
	if err == nil {
		t.Fatal("expected error when validating expired token, got nil")
	}

	// Test 3: Invalid Secret Key Validation
	wrongSecret := "wrongsecretkey123456789012345678"
	_, err = ValidateToken(token, wrongSecret)
	if err == nil {
		t.Fatal("expected error when validating token with wrong secret, got nil")
	}
}
