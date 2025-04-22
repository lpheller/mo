package utils

import (
	"fmt"
	"os"
	"strings"
)

type EnvManager struct {
	Path string
}

func NewEnvManager(path string) *EnvManager {
	return &EnvManager{Path: path}
}

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

func EnsureRequiredEnvVars(context string) (map[string]string, error) {
	var requiredKeys []string

	if context == "push" {
		requiredKeys = []string{"PUSH_SSH_USER", "PUSH_HOST", "PUSH_PROJECT_DIR"}
	} else if context == "pull" {
		requiredKeys = []string{"PULL_SSH_USER", "PULL_HOST", "PULL_PROJECT_DIR"}
	} else {
		return nil, fmt.Errorf("invalid context: %s", context)
	}

	envManager := NewEnvManager(".env")
	envVars := make(map[string]string)

	for _, key := range requiredKeys {
		value, found, err := envManager.GetVar(key)
		if err != nil {
			return nil, fmt.Errorf("error reading key %s from .env: %v", key, err)
		}
		if !found {
			value, err = promptForEnvVar(key)
			if err != nil {
				return nil, err
			}
			if err := envManager.SetVar(key, value); err != nil {
				return nil, fmt.Errorf("error saving key %s to .env: %v", key, err)
			}
		}
		envVars[key] = value
	}

	return envVars, nil
}

func promptForEnvVar(key string) (string, error) {
	fmt.Printf("The required environment variable '%s' is missing. Please enter its value: ", key)
	var value string
	_, err := fmt.Scanln(&value)
	if err != nil {
		return "", fmt.Errorf("error reading input for %s: %v", key, err)
	}
	return value, nil
}
