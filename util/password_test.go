package util

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := "secret123"

	hashedPassword1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if len(hashedPassword1) == 0 {
		t.Fatal("expected non-empty hashed password")
	}

	err = CheckPassword(password, hashedPassword1)
	if err != nil {
		t.Fatalf("expected password match, got error: %v", err)
	}

	wrongPassword := "wrongsecret"
	err = CheckPassword(wrongPassword, hashedPassword1)
	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Fatalf("expected ErrMismatchedHashAndPassword, got: %v", err)
	}

	hashedPassword2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password second time: %v", err)
	}
	if hashedPassword1 == hashedPassword2 {
		t.Fatal("expected different hashes due to bcrypt salt generation")
	}
}
