package utils

import (
	"fmt"
	"os"
	"strings"
)

type EnvManager struct {
	Path string
}

// New creates a new instance with the given .env file path
func NewEnvManager(path string) *EnvManager {
	return &EnvManager{Path: path}
}

// GetVar returns the value of a given variable, and a boolean indicating if it was found
func (e *EnvManager) GetVar(key string) (string, bool, error) {
	data, err := os.ReadFile(e.Path)
	if err != nil {
		return "", false, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1]), true, nil
		}
	}
	return "", false, nil
}

func (e *EnvManager) SetVar(key, value string) error {
	data, err := os.ReadFile(e.Path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist — create it with the variable
			line := fmt.Sprintf("%s=%s\n", key, value)
			return os.WriteFile(e.Path, []byte(line), 0644)
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Check both commented and uncommented lines
		isCommented := strings.HasPrefix(trimmed, "#")
		content := trimmed
		if isCommented {
			content = strings.TrimSpace(trimmed[1:])
		}

		parts := strings.SplitN(content, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			lines[i] = fmt.Sprintf("%s=%s", key, value) // overwrite and uncomment
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	content := strings.Join(lines, "\n")
	return os.WriteFile(e.Path, []byte(content), 0644)
}
