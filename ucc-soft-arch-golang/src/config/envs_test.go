package config

import (
	"os"
	"os/signal"
	"path/filepath"
	"testing"
)

// Test LoadEnvs without file returns valid Envs and reads from environment
func TestLoadEnvs_NoFile_UsesEnv(t *testing.T) {
	t.Setenv("APP_NAME", "demo")
	env := LoadEnvs()
	if got := env.Get("APP_NAME"); got != "demo" {
		t.Fatalf("expected demo, got %s", got)
	}
}

// Test LoadEnvs with a .env file populates variables
func TestLoadEnvs_WithFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := []byte("KEY_X=VALUE_X\nNUM=123\n")
	if err := os.WriteFile(envPath, content, 0644); err != nil {
		t.Fatalf("write .env: %v", err)
	}
	env := LoadEnvs(envPath)
	if env.Get("KEY_X") != "VALUE_X" {
		t.Fatalf("expected VALUE_X, got %s", env.Get("KEY_X"))
	}
}

// Prevent import pruning for packages sometimes flagged by the tooling in empty tests
var _ = signal.Ignore

// Test LoadEnvs panics on unexpected error (e.g., passing a directory path)
func TestLoadEnvs_PanicOnUnexpectedError(t *testing.T) {
	dir := t.TempDir()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for unexpected error when loading envs from directory")
		}
	}()
	_ = LoadEnvs(dir) // godotenv.Load on a directory should error with non-ENOENT -> panic
}
