package commands

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func EnvMailDev(cliContext *cli.Context) error {

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
			line = "MAIL_HOST=127.0.0.1"
		} else if strings.HasPrefix(line, "MAIL_PORT=") {
			line = "MAIL_PORT=2525"
		} else if strings.HasPrefix(line, "MAIL_USERNAME=") {
			line = "MAIL_USERNAME="
		} else if strings.HasPrefix(line, "MAIL_PASSWORD=") {
			line = "MAIL_PASSWORD="
		} else if strings.HasPrefix(line, "MAIL_ENCRYPTION=") {
			line = "MAIL_ENCRYPTION="
		} else if strings.HasPrefix(line, "MAIL_FROM_ADDRESS=") {
			line = "MAIL_FROM_ADDRESS="
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
