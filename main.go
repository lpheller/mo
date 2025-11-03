package main

import (
	"log"
	"os"

	"mo/commands"
	"mo/config"

	"github.com/urfave/cli/v2"
)

func main() {
	if _, err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "db",
				Usage: "Database management",
				Subcommands: []*cli.Command{
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "Create a new database",
						Action:  commands.CreateDatabase,
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "List all databases",
						Action:  commands.ListDatabases,
					},
					{
						Name:    "open",
						Aliases: []string{"o"},
						Usage:   "Open the database in the default editor",
						Action:  commands.OpenDatabase,
					},
				},
			},
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
				Name:  "env",
				Usage: "Environment management",
				Subcommands: []*cli.Command{
					{
						Name:   "sqlite",
						Usage:  "Set the DB_CONNECTION to sqlite",
						Action: commands.EnvSqlite,
					},
					{
						Name:   "mailtrap",
						Usage:  "Set the mail driver to mailtrap",
						Action: commands.EnvMailtrap,
					},
					{
						Name:   "maildev",
						Usage:  "Set the mail driver to mail-dev",
						Action: commands.EnvMailDev,
					},
					{
						Name:   "sync",
						Usage:  "Sync the .env file with .env.example",
						Action: commands.SyncEnv,
					},
				},
			},
			{
				Name:   "env:sqlite",
				Usage:  "Set the DB_CONNECTION to sqlite",
				Action: commands.EnvSqlite,
			},
			{
				Name:   "env:mailtrap",
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
				Aliases:     []string{"sync:env"},
				Usage:       "Sync the .env file with .env.example",
				Description: `Sync the .env file with .env.example`,
				Action:      commands.SyncEnv,
			},
			{
				Name:    "config",
				Aliases: []string{"cfg", "qc"},
				Usage:   "Quickly open configuration files",
				Action:  commands.QuickConfig,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "editor",
						Usage:   "Specify a custom editor",
						Aliases: []string{"e"},
					},
				},
			},
			{
				Name:    "setup",
				Aliases: []string{"s"},
				Usage:   "Setup a project by running appropriate commands",
				Action:  commands.CheckProject,
			},
			{
				Name:   "pull",
				Usage:  "Pull storage or database from a remote server",
				Action: commands.Pull,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "storage",
						Usage: "Pull the storage folder",
					},
					&cli.BoolFlag{
						Name:  "database",
						Usage: "Pull the database",
					},
				},
			},
			{
				Name:   "push",
				Usage:  "Push storage or database to a remote server",
				Action: commands.Push,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "storage",
						Usage: "Push the storage folder",
					},
					&cli.BoolFlag{
						Name:  "database",
						Usage: "Push the database",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
