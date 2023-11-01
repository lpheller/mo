package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUser     string `json:"DB_USERNAME"`
	DbPassword string `json:"DB_PASSWORD"`
	DbHost     string `json:"DB_HOST"`
	DbPort     string `json:"DB_PORT"`
	Editor     string `json:"EDITOR"`

	MailtrapUsername string `json:"MAILTRAP_USERNAME"`
	MailtrapPassword string `json:"MAILTRAP_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "mortimer", "config.json")

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			return nil, err
		}

		f, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		config := &Config{
			DBUser:     "root",
			DbPassword: "",
			DbHost:     "127.0.0.1",
			DbPort:     "3306",
			Editor:     "vscode",

			MailtrapUsername: "",
			MailtrapPassword: "",
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, err
		}

		f.Write(configBytes)
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
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, ".config", "mortimer", "config.json")

	return configPath, nil
}