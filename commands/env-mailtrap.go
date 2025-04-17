package commands

import (
	"bufio"
	"fmt"
	"mo/config"
	"os"
	"strings"

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

	// Open the .env file
	file, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("failed to open .env file: %v", err)
	}
	defer file.Close()

	// Create a map for the replacements
	replacements := map[string]string{
		"MAIL_MAILER=":       "MAIL_MAILER=smtp",
		"MAIL_HOST=":         "MAIL_HOST=smtp.mailtrap.io",
		"MAIL_PORT=":         "MAIL_PORT=2525",
		"MAIL_USERNAME=":     "MAIL_USERNAME=" + mailtrapUsername,
		"MAIL_PASSWORD=":     "MAIL_PASSWORD=" + mailtrapPassword,
		"MAIL_ENCRYPTION=":   "MAIL_ENCRYPTION=tls",
		"MAIL_FROM_ADDRESS=": "MAIL_FROM_ADDRESS=mail@project.test",
	}

	var lines []string
	scanner := bufio.NewScanner(file)

	// Read each line and apply necessary replacements
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line matches any key from the replacement map
		for prefix, replacement := range replacements {
			if strings.HasPrefix(line, prefix) {
				line = replacement
				break
			}
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}

	// Re-open the .env file in write mode to save the updated lines
	file, err = os.Create(".env")
	if err != nil {
		return fmt.Errorf("failed to create .env file: %v", err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := w.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to .env file: %v", err)
		}
	}

	// Ensure everything is written to disk
	if err := w.Flush(); err != nil {
		return fmt.Errorf("error flushing to .env file: %v", err)
	}

	fmt.Println(".env file updated with Mailtrap credentials.")
	return nil
}
