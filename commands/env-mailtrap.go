package commands

import (
	"bufio"
	"log"
	"mo/config"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func EnvMailtrap(c *cli.Context) error {

	cfg, err := config.LoadConfig()

	mailtrapUsername := cfg.MailtrapUsername
	mailtrapPassword := cfg.MailtrapPassword

	if mailtrapUsername == "" || mailtrapPassword == "" {
		log.Println("Your mailtrap credentials are missing. Please run \"mo config:edit\" to set them.")
		return nil
	}

	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "MAIL_MAILER=") {
			line = "MAIL_MAILER=smtp"
		} else if strings.HasPrefix(line, "MAIL_HOST=") {
			line = "MAIL_HOST=smtp.mailtrap.io"
		} else if strings.HasPrefix(line, "MAIL_PORT=") {
			line = "MAIL_PORT=2525"
		} else if strings.HasPrefix(line, "MAIL_USERNAME=") {
			line = "MAIL_USERNAME=" + mailtrapUsername
		} else if strings.HasPrefix(line, "MAIL_PASSWORD=") {
			line = "MAIL_PASSWORD=" + mailtrapPassword
		} else if strings.HasPrefix(line, "MAIL_ENCRYPTION=") {
			line = "MAIL_ENCRYPTION=tls"
		} else if strings.HasPrefix(line, "MAIL_FROM_ADDRESS=") {
			line = "MAIL_FROM_ADDRESS=mail@project.test"
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	file, err = os.Create(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := w.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()

	return nil

}
