package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func OpenDatabase(cliContext *cli.Context) error {
	_, err := os.Stat(".env")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	envExists := !os.IsNotExist(err)

	if !envExists {
		return fmt.Errorf(".env file not found")
	}

	envContent, err := os.ReadFile(".env")
	if err != nil {
		return err
	}

	if strings.Contains(string(envContent), "DB_CONNECTION=sqlite") {
		if _, err := os.Stat("database/database.sqlite"); !os.IsNotExist(err) {
			return exec.Command("open", "database/database.sqlite").Run()
		}
	}

	var dbConnection, dbUsername, dbPassword, dbHost, dbPort, dbDatabase string

	lines := strings.Split(string(envContent), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "DB_CONNECTION":
				dbConnection = value
			case "DB_USERNAME":
				dbUsername = value
			case "DB_PASSWORD":
				dbPassword = value
			case "DB_HOST":
				dbHost = value
			case "DB_PORT":
				dbPort = value
			case "DB_DATABASE":
				dbDatabase = value
			}
		}
	}

	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s", dbConnection, dbUsername, dbPassword, dbHost, dbPort, dbDatabase)

	return exec.Command("open", connStr).Run()
}
