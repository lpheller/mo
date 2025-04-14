package utils

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// LoadEnv lädt eine .env-Datei und validiert erforderliche Schlüssel.
func LoadEnv(filePath string, requiredKeys []string) (map[string]string, error) {
    env := make(map[string]string)

    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("error opening env file: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
            env[key] = value
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading env file: %v", err)
    }

    for _, key := range requiredKeys {
        if _, exists := env[key]; !exists {
            return nil, fmt.Errorf("missing required key '%s' in env file", key)
        }
    }

    return env, nil
}