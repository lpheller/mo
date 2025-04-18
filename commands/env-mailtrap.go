package commands

import (
	"fmt"
	"mo/config"
	"mo/utils"

	"github.com/urfave/cli/v2"
)

func EnvMailtrap(cliContext *cli.Context) error {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	mailtrapUsername := cfg.MailtrapUsername
	mailtrapPassword := cfg.MailtrapPassword

	if mailtrapUsername == "" || mailtrapPassword == "" {
		return fmt.Errorf("your Mailtrap credentials are missing. Please run \"mo config:edit\" to set them")
	}

	// Initialize EnvManager
	envManager := utils.NewEnvManager(".env")

	// Define the replacements
	replacements := map[string]string{
		"MAIL_MAILER":       "smtp",
		"MAIL_HOST":         "smtp.mailtrap.io",
		"MAIL_PORT":         "2525",
		"MAIL_USERNAME":     mailtrapUsername,
		"MAIL_PASSWORD":     mailtrapPassword,
		"MAIL_ENCRYPTION":   "tls",
		"MAIL_FROM_ADDRESS": "mail@project.test",
	}

	// Apply the replacements using EnvManager
	for key, value := range replacements {
		if err := envManager.SetVar(key, value); err != nil {
			return fmt.Errorf("failed to set %s in .env file: %v", key, err)
		}
	}

	fmt.Println(".env file updated with Mailtrap credentials.")
	return nil
}
