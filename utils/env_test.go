package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvManager_GetVar(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	content := `DB_HOST=localhost
DB_PORT=3306
# COMMENTED_VAR=value
EMPTY_VAR=
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	manager := NewEnvManager(envPath)

	tests := []struct {
		name      string
		key       string
		wantValue string
		wantFound bool
	}{
		{"existing variable", "DB_HOST", "localhost", true},
		{"existing port", "DB_PORT", "3306", true},
		{"commented variable", "COMMENTED_VAR", "", false},
		{"nonexistent variable", "NONEXISTENT", "", false},
		{"empty variable", "EMPTY_VAR", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found, err := manager.GetVar(tt.key)
			if err != nil {
				t.Fatalf("GetVar() error = %v", err)
			}
			if found != tt.wantFound {
				t.Errorf("GetVar() found = %v, want %v", found, tt.wantFound)
			}
			if value != tt.wantValue {
				t.Errorf("GetVar() value = %v, want %v", value, tt.wantValue)
			}
		})
	}
}

func TestEnvManager_SetVar(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	content := `DB_HOST=localhost
# DB_PORT=3306
DB_USER=root
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	manager := NewEnvManager(envPath)

	tests := []struct {
		name      string
		key       string
		value     string
		wantValue string
	}{
		{"update existing", "DB_HOST", "127.0.0.1", "127.0.0.1"},
		{"uncomment and set", "DB_PORT", "5432", "5432"},
		{"add new", "NEW_VAR", "newvalue", "newvalue"},
		{"update with special chars", "DB_USER", "admin@localhost", "admin@localhost"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := manager.SetVar(tt.key, tt.value); err != nil {
				t.Fatalf("SetVar() error = %v", err)
			}

			value, found, err := manager.GetVar(tt.key)
			if err != nil {
				t.Fatalf("GetVar() error = %v", err)
			}
			if !found {
				t.Errorf("Variable %s not found after SetVar()", tt.key)
			}
			if value != tt.wantValue {
				t.Errorf("GetVar() value = %v, want %v", value, tt.wantValue)
			}
		})
	}
}

func TestEnvManager_SetVar_NewFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env.new")

	manager := NewEnvManager(envPath)

	if err := manager.SetVar("NEW_VAR", "value"); err != nil {
		t.Fatalf("SetVar() error = %v", err)
	}

	value, found, err := manager.GetVar("NEW_VAR")
	if err != nil {
		t.Fatalf("GetVar() error = %v", err)
	}
	if !found {
		t.Error("Variable not found after creating new file")
	}
	if value != "value" {
		t.Errorf("GetVar() value = %v, want %v", value, "value")
	}

	// Verify file was created
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Error("File was not created")
	}
}

func TestEnvManager_GetVar_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env.notexist")

	manager := NewEnvManager(envPath)

	_, _, err := manager.GetVar("SOME_VAR")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
