package jwt

import (
	"testing"

	"github.com/google/uuid"
)

func TestSignAndVerifyToken_Success(t *testing.T) {
	id := uuid.New()
	role := 2
	token := SignDocument(id, role)
	if token == "" {
		t.Fatalf("expected non-empty token")
	}
	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("unexpected error verifying token: %v", err)
	}
	if claims["id"] == nil || claims["role"] == nil {
		t.Fatalf("expected id and role claims present")
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	_, err := VerifyToken("invalid.token.format")
	if err == nil {
		t.Fatalf("expected error for invalid token")
	}
}
