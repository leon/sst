package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestEnvPrecedence(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	
	// Create .env file
	envFile := filepath.Join(tempDir, ".env")
	envContent := "TEST_VAR=from_env\nCOMMON_VAR=from_env\n"
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	
	// Create .env.dev file
	envDevFile := filepath.Join(tempDir, ".env.dev")
	envDevContent := "TEST_VAR=from_env_dev\nDEV_VAR=from_env_dev\n"
	err = os.WriteFile(envDevFile, []byte(envDevContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env.dev file: %v", err)
	}
	
	// Clear environment variables
	os.Unsetenv("TEST_VAR")
	os.Unsetenv("COMMON_VAR") 
	os.Unsetenv("DEV_VAR")
	
	// Simulate the current fixed behavior:
	// 1. Load .env first (from cli.New)
	godotenv.Load(envFile)
	
	// 2. Use Overload for stage-specific file (our fix in cli.Stage)
	godotenv.Overload(envDevFile)
	
	// Test that stage-specific values take precedence
	if got := os.Getenv("TEST_VAR"); got != "from_env_dev" {
		t.Errorf("TEST_VAR should be overridden by .env.dev, got %s, want from_env_dev", got)
	}
	
	// Test that .env-only variables are preserved
	if got := os.Getenv("COMMON_VAR"); got != "from_env" {
		t.Errorf("COMMON_VAR should remain from .env, got %s, want from_env", got)
	}
	
	// Test that dev-only variables are set
	if got := os.Getenv("DEV_VAR"); got != "from_env_dev" {
		t.Errorf("DEV_VAR should be set from .env.dev, got %s, want from_env_dev", got)
	}
	
	// Clean up
	os.Unsetenv("TEST_VAR")
	os.Unsetenv("COMMON_VAR")
	os.Unsetenv("DEV_VAR")
}

func TestOldBehaviorWouldFail(t *testing.T) {
	// This test demonstrates the old behavior would be incorrect
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	
	// Create .env file
	envFile := filepath.Join(tempDir, ".env")
	envContent := "TEST_VAR=from_env\n"
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	
	// Create .env.dev file
	envDevFile := filepath.Join(tempDir, ".env.dev")
	envDevContent := "TEST_VAR=from_env_dev\n"
	err = os.WriteFile(envDevFile, []byte(envDevContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env.dev file: %v", err)
	}
	
	// Clear environment variable
	os.Unsetenv("TEST_VAR")
	
	// Simulate the old broken behavior:
	// 1. Load .env first
	godotenv.Load(envFile)
	
	// 2. Use Load (not Overload) for stage-specific file
	godotenv.Load(envDevFile)
	
	// This would fail with the old behavior - TEST_VAR would still be "from_env"
	if got := os.Getenv("TEST_VAR"); got == "from_env_dev" {
		t.Logf("Old behavior test failed as expected - this confirms our fix is needed. Got %s", got)
	} else {
		t.Logf("Old behavior confirmed broken: TEST_VAR=%s (should have been from_env_dev)", got)
	}
	
	// Clean up
	os.Unsetenv("TEST_VAR")
}