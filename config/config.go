package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUser           string `json:"DB_USERNAME"`
	DbPassword       string `json:"DB_PASSWORD"`
	DbHost           string `json:"DB_HOST"`
	DbPort           string `json:"DB_PORT"`
	Editor           string `json:"EDITOR"`
	MailtrapUsername string `json:"MAILTRAP_USERNAME"`
	MailtrapPassword string `json:"MAILTRAP_PASSWORD"`
	ConfigPaths      map[string]string
}

// DefaultConfig creates a default configuration.
func DefaultConfig() *Config {
	return &Config{
		DBUser:           "root",
		DbPassword:       "",
		DbHost:           "127.0.0.1",
		DbPort:           "3306",
		Editor:           "vscode",
		MailtrapUsername: "",
		MailtrapPassword: "",
		ConfigPaths:      map[string]string{},
	}
}

// LoadConfig loads the configuration from a JSON file or creates it if it doesn't exist.
func LoadConfig() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		err = createDefaultConfig(configPath)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ConfigPath returns the path to the configuration file.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "mortimer", "config.json"), nil
}

func createDefaultConfig(configPath string) error {
	config := DefaultConfig()
	configBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(configBytes)
	return err
}
