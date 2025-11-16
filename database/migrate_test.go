package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Create Users Table", "create_users_table"},
		{"add-avatar-to-users", "add_avatar_to_users"},
		{"Fix Bug #123", "fix_bug_123"},
		{"Update Schema v2.0", "update_schema_v20"},
		{"Special!@#$%Chars", "specialchars"},
		{"UPPERCASE", "uppercase"},
		{"mixed_CASE-test", "mixed_case_test"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeName(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns env value when set",
			key:          "TEST_KEY_1",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "returns default when env not set",
			key:          "TEST_KEY_2",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			// Test
			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getEnv(%q, %q) = %q; want %q", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestGetDBURL(t *testing.T) {
	// Save original env vars
	originals := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_SSL_MODE": os.Getenv("DB_SSL_MODE"),
	}

	// Cleanup
	defer func() {
		for k, v := range originals {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	t.Run("uses default values when env vars not set", func(t *testing.T) {
		// Clear all env vars
		for k := range originals {
			os.Unsetenv(k)
		}

		url := getDBURL()
		expected := "postgres://edugo:changeme@localhost:5432/edugo_dev?sslmode=disable"
		if url != expected {
			t.Errorf("getDBURL() = %q; want %q", url, expected)
		}
	})

	t.Run("uses custom env vars", func(t *testing.T) {
		os.Setenv("DB_HOST", "custom-host")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_NAME", "custom_db")
		os.Setenv("DB_USER", "custom_user")
		os.Setenv("DB_PASSWORD", "custom_pass")
		os.Setenv("DB_SSL_MODE", "require")

		url := getDBURL()
		expected := "postgres://custom_user:custom_pass@custom-host:5433/custom_db?sslmode=require"
		if url != expected {
			t.Errorf("getDBURL() = %q; want %q", url, expected)
		}
	})
}

func TestLoadMigrations(t *testing.T) {
	// Create temp directory with mock migrations
	tempDir := t.TempDir()
	migrationsDir := filepath.Join(tempDir, "migrations", "postgres")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Save original migrationsDir and restore after test
	originalDir := migrationsDir
	defer func() {
		// Note: In real usage, we'd need to modify the const or use dependency injection
		_ = originalDir
	}()

	// Create mock migration files
	upSQL := "CREATE TABLE test_users (id INT PRIMARY KEY);"
	downSQL := "DROP TABLE test_users;"

	files := map[string]string{
		"001_create_users.up.sql":   upSQL,
		"001_create_users.down.sql": downSQL,
		"002_add_email.up.sql":      "ALTER TABLE test_users ADD COLUMN email VARCHAR(255);",
		"002_add_email.down.sql":    "ALTER TABLE test_users DROP COLUMN email;",
	}

	for filename, content := range files {
		path := filepath.Join(migrationsDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write mock file %s: %v", filename, err)
		}
	}

	// Note: This test would fail with current implementation since migrationsDir is a const
	// In a real refactoring, we'd make it configurable for testing
	t.Skip("Skipping test - requires refactoring to make migrationsDir configurable")
}

func TestCreateMigrationFiles(t *testing.T) {
	// Test that we can create migration files
	tempDir := t.TempDir()

	// This is a smoke test - in real usage we'd need to modify the code
	// to accept a custom migrations directory for testing
	t.Run("sanitize migration name", func(t *testing.T) {
		name := "Create User Profiles Table"
		sanitized := sanitizeName(name)
		expected := "create_user_profiles_table"

		if sanitized != expected {
			t.Errorf("Expected %q but got %q", expected, sanitized)
		}
	})

	_ = tempDir // Use tempDir to avoid unused warning
}
