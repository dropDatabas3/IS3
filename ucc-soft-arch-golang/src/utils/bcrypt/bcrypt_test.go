package bcrypt

import (
	"testing"
)

func TestHashAndCompare_Success(t *testing.T) {
	plain := "S3cret!"
	hashed, err := HasPassword(plain)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}
	if hashed == plain || len(hashed) == 0 {
		t.Fatalf("hashed password should not equal plain text or be empty")
	}
	if ok := ComparePassword(plain, hashed); !ok {
		t.Fatalf("expected ComparePassword to return true for correct password")
	}
}

func TestHashAndCompare_Failure(t *testing.T) {
	plain := "S3cret!"
	hashed, err := HasPassword(plain)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}
	if ok := ComparePassword("wrong", hashed); ok {
		t.Fatalf("expected ComparePassword to return false for wrong password")
	}
}
