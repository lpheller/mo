package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"DBUser", cfg.DBUser, "root"},
		{"DBHost", cfg.DBHost, "127.0.0.1"},
		{"DBPort", cfg.DBPort, "3306"},
		{"Editor", cfg.Editor, "vscode"},
		{"DBPassword", cfg.DBPassword, ""},
		{"MailtrapUsername", cfg.MailtrapUsername, ""},
		{"MailtrapPassword", cfg.MailtrapPassword, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	if cfg.ConfigPaths == nil {
		t.Error("ConfigPaths should be initialized, got nil")
	}
}

func TestLoadConfig_CreatesDefault(t *testing.T) {
	tmpDir := t.TempDir()
	tmpHome := filepath.Join(tmpDir, "home")

	// Override configPathFunc for testing
	originalConfigPath := configPathFunc
	configPathFunc = func() (string, error) {
		return filepath.Join(tmpHome, ".config", "mortimer", "config.json"), nil
	}
	defer func() { configPathFunc = originalConfigPath }()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DBUser != "root" {
		t.Errorf("Default config not created correctly, DBUser = %v, want root", cfg.DBUser)
	}

	if cfg.Editor != "vscode" {
		t.Errorf("Default config not created correctly, Editor = %v, want vscode", cfg.Editor)
	}

	// Verify file was actually created
	configPath, _ := ConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestLoadConfig_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	testConfig := &Config{
		DBUser:           "testuser",
		DBPassword:       "testpass",
		DBHost:           "testhost",
		DBPort:           "5432",
		Editor:           "vim",
		MailtrapUsername: "mailtrap@test.com",
		MailtrapPassword: "mailtrappass",
		ConfigPaths:      map[string]string{"test": "/path/to/test"},
	}

	configBytes, _ := json.MarshalIndent(testConfig, "", "  ")
	if err := os.WriteFile(configPath, configBytes, 0644); err != nil {
		t.Fatal(err)
	}

	// Override configPathFunc for testing
	originalConfigPath := configPathFunc
	configPathFunc = func() (string, error) {
		return configPath, nil
	}
	defer func() { configPathFunc = originalConfigPath }()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"DBUser", cfg.DBUser, "testuser"},
		{"DBPassword", cfg.DBPassword, "testpass"},
		{"DBHost", cfg.DBHost, "testhost"},
		{"DBPort", cfg.DBPort, "5432"},
		{"Editor", cfg.Editor, "vim"},
		{"MailtrapUsername", cfg.MailtrapUsername, "mailtrap@test.com"},
		{"MailtrapPassword", cfg.MailtrapPassword, "mailtrappass"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	if cfg.ConfigPaths["test"] != "/path/to/test" {
		t.Errorf("ConfigPaths not loaded correctly")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Write invalid JSON
	if err := os.WriteFile(configPath, []byte("invalid json {{{"), 0644); err != nil {
		t.Fatal(err)
	}

	// Override configPathFunc for testing
	originalConfigPath := configPathFunc
	configPathFunc = func() (string, error) {
		return configPath, nil
	}
	defer func() { configPathFunc = originalConfigPath }()

	_, err := LoadConfig()
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestConfigPath(t *testing.T) {
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() returned empty string")
	}

	// Check that path contains expected components
	if !filepath.IsAbs(path) {
		t.Errorf("ConfigPath() should return absolute path, got %v", path)
	}

	if filepath.Base(path) != "config.json" {
		t.Errorf("ConfigPath() should end with config.json, got %v", path)
	}
}
