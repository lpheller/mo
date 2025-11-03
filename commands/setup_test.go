package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")

	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"existing file", existingFile, true},
		{"non-existing file", filepath.Join(tmpDir, "notexists.txt"), false},
		{"directory", tmpDir, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.path); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")

	content := []byte("test content for copy")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}

	destContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(destContent) != string(content) {
		t.Errorf("File content = %v, want %v", string(destContent), string(content))
	}
}

func TestCopyFile_NonExistentSource(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "nonexistent.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")

	err := copyFile(srcFile, dstFile)
	if err == nil {
		t.Error("copyFile() expected error for non-existent source, got nil")
	}
}

func TestHasEnvAppKey(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			"with key",
			"APP_KEY=base64:test123\nDB_HOST=localhost",
			true,
		},
		{
			"without key",
			"DB_HOST=localhost\nDB_PORT=3306",
			false,
		},
		{
			"empty key",
			"APP_KEY=\nDB_HOST=localhost",
			false,
		},
		{
			"key with whitespace",
			"APP_KEY=base64:test123  \nDB_HOST=localhost",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			envFile := filepath.Join(tmpDir, ".env")
			if err := os.WriteFile(envFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			// Change to temp directory
			oldWd, _ := os.Getwd()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldWd)

			got, err := hasEnvAppKey()
			if err != nil {
				t.Fatalf("hasEnvAppKey() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("hasEnvAppKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasEnvAppKey_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	_, err := hasEnvAppKey()
	if err == nil {
		t.Error("hasEnvAppKey() expected error for missing file, got nil")
	}
}

func TestHasNpmScript(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		scriptName string
		want       bool
	}{
		{
			"has build script",
			`{
  "name": "test",
  "scripts": {
    "build": "vite build",
    "dev": "vite"
  },
}`,
			"build",
			true,
		},
		{
			"no build script",
			`{
  "name": "test",
  "scripts": {
    "dev": "vite"
  },
}`,
			"build",
			false,
		},
		{
			"no scripts section",
			`{
  "name": "test",
  "version": "1.0.0"
}`,
			"build",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			packageFile := filepath.Join(tmpDir, "package.json")
			if err := os.WriteFile(packageFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			// Change to temp directory
			oldWd, _ := os.Getwd()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldWd)

			got, err := hasNpmScript(tt.scriptName)
			if err != nil {
				t.Fatalf("hasNpmScript() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("hasNpmScript() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnsureEnvFile(t *testing.T) {
	t.Run("env file exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		envFile := filepath.Join(tmpDir, ".env")
		if err := os.WriteFile(envFile, []byte("DB_HOST=localhost"), 0644); err != nil {
			t.Fatal(err)
		}

		oldWd, _ := os.Getwd()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldWd)

		if err := ensureEnvFile(); err != nil {
			t.Errorf("ensureEnvFile() error = %v", err)
		}
	})

	t.Run("copy from example", func(t *testing.T) {
		tmpDir := t.TempDir()
		exampleFile := filepath.Join(tmpDir, ".env.example")
		envFile := filepath.Join(tmpDir, ".env")
		exampleContent := []byte("DB_HOST=localhost\nDB_PORT=3306")

		if err := os.WriteFile(exampleFile, exampleContent, 0644); err != nil {
			t.Fatal(err)
		}

		oldWd, _ := os.Getwd()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldWd)

		if err := ensureEnvFile(); err != nil {
			t.Fatalf("ensureEnvFile() error = %v", err)
		}

		content, err := os.ReadFile(envFile)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != string(exampleContent) {
			t.Errorf("Content mismatch, got %v, want %v", string(content), string(exampleContent))
		}
	})

	t.Run("no example file", func(t *testing.T) {
		tmpDir := t.TempDir()

		oldWd, _ := os.Getwd()
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldWd)

		err := ensureEnvFile()
		if err == nil {
			t.Error("ensureEnvFile() expected error when .env.example missing, got nil")
		}
	})
}
