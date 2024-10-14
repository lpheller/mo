package main

import (
	"log"
	"os"

	"mo/commands"
	"mo/config"

	"github.com/urfave/cli/v2"
)

func main() {
	_, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "db:open",
				Aliases: []string{"opendb"},
				Usage:   "Open the database in the default editor",
				Action:  commands.OpenDatabase,
			},
			{
				Name:    "db:list",
				Aliases: []string{"listdb"},
				Usage:   "List all databases",
				Action:  commands.ListDatabases,
			},
			{
				Name:    "db:create",
				Aliases: []string{"createdb"},
				Usage:   "Create a new database",
				Action:  commands.CreateDatabase,
			},
			{
				Name: "env:sqlite",
				// Aliases: []string{"envsqlite"},
				Usage:  "Set the DB_CONNECTION to sqlite",
				Action: commands.EnvSqlite,
			},
			{
				Name: "env:mailtrap",
				// Aliases: []string{"envmailtrap"},
				Usage:  "Set the mail driver to mailtrap",
				Action: commands.EnvMailtrap,
			},
			{
				Name:   "env:maildev",
				Usage:  "Set the mail driver to mail-dev",
				Action: commands.EnvMailDev,
			},
			{
				Name:    "config:edit",
				Aliases: []string{"edit:config"},
				Usage:   "Edit the Mortimer config file",
				Action:  commands.EditConfig,
			},
			{
				Name:        "env:sync",
				Aliases:     []string{"snyc:env"},
				Usage:       "Sync the .env file with .env.example",
				Description: `Sync the .env file with .env.example`,
				Action:      commands.SyncEnv,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
