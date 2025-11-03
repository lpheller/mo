package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUser           string            `json:"db_user"`
	DBPassword       string            `json:"db_password"`
	DBHost           string            `json:"db_host"`
	DBPort           string            `json:"db_port"`
	Editor           string            `json:"editor"`
	MailtrapUsername string            `json:"mailtrap_username"`
	MailtrapPassword string            `json:"mailtrap_password"`
	ConfigPaths      map[string]string `json:"config_paths"`
}

func DefaultConfig() *Config {
	return &Config{
		DBUser:           "root",
		DBPassword:       "",
		DBHost:           "127.0.0.1",
		DBPort:           "3306",
		Editor:           "vscode",
		MailtrapUsername: "",
		MailtrapPassword: "",
		ConfigPaths:      map[string]string{},
	}
}

// configPathFunc is a variable that holds the function to get the config path
// This allows for easier testing by allowing the function to be overridden
var configPathFunc = defaultConfigPath

func LoadConfig() (*Config, error) {
	configPath, err := configPathFunc()
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

func ConfigPath() (string, error) {
	return configPathFunc()
}

func defaultConfigPath() (string, error) {
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
