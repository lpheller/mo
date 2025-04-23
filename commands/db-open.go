package commands

import (
	"fmt"
	"os"
	"os/exec"

	"mo/utils"

	"github.com/urfave/cli/v2"
)

func OpenDatabase(cliContext *cli.Context) error {
	envManager := utils.NewEnvManager(".env")

	if _, err := os.Stat(envManager.Path); os.IsNotExist(err) {
		return fmt.Errorf(".env file not found")
	}

	dbConnection, found, err := envManager.GetVar("DB_CONNECTION")
	if err != nil {
		return fmt.Errorf("error reading DB_CONNECTION: %v", err)
	}
	if !found {
		return fmt.Errorf("DB_CONNECTION not found in .env")
	}

	if dbConnection == "sqlite" {
		if _, err := os.Stat("database/database.sqlite"); !os.IsNotExist(err) {
			return exec.Command("open", "database/database.sqlite").Run()
		}
		return fmt.Errorf("SQLite database file not found")
	}

	dbUsername, _, _ := envManager.GetVar("DB_USERNAME")
	dbPassword, _, _ := envManager.GetVar("DB_PASSWORD")
	dbHost, _, _ := envManager.GetVar("DB_HOST")
	dbPort, _, _ := envManager.GetVar("DB_PORT")
	dbDatabase, _, _ := envManager.GetVar("DB_DATABASE")

	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s", dbConnection, dbUsername, dbPassword, dbHost, dbPort, dbDatabase)

	return exec.Command("open", connStr).Run()
}
