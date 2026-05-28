package common

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHash_NonEmpty(t *testing.T) {
	hash := BcryptHash("password123")
	if hash == "" {
		t.Fatal("BcryptHash returned empty string")
	}
}

func TestBcryptHash_VerifiesAgainstSource(t *testing.T) {
	const password = "s3cr3t-pass"
	hash := BcryptHash(password)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		t.Fatalf("bcrypt.CompareHashAndPassword failed: %v", err)
	}
}

func TestBcryptHash_DifferentSaltEachCall(t *testing.T) {
	const password = "same-password"
	h1 := BcryptHash(password)
	h2 := BcryptHash(password)
	if h1 == h2 {
		t.Fatalf("expected different salts/hashes, got identical: %s", h1)
	}
}
