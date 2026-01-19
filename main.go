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

	// Define commands once to avoid duplication
	dbCreateCmd := &cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "Create a new database",
		Action:  commands.CreateDatabase,
	}

	dbListCmd := &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all databases",
		Action:  commands.ListDatabases,
	}

	dbOpenCmd := &cli.Command{
		Name:    "open",
		Aliases: []string{"o"},
		Usage:   "Open the database in the default editor",
		Action:  commands.OpenDatabase,
	}

	envSqliteCmd := &cli.Command{
		Name:   "sqlite",
		Usage:  "Set the DB_CONNECTION to sqlite",
		Action: commands.EnvSqlite,
	}

	envMailtrapCmd := &cli.Command{
		Name:   "mailtrap",
		Usage:  "Set the mail driver to mailtrap",
		Action: commands.EnvMailtrap,
	}

	envMaildevCmd := &cli.Command{
		Name:   "maildev",
		Usage:  "Set the mail driver to mail-dev",
		Action: commands.EnvMailDev,
	}

	envSyncCmd := &cli.Command{
		Name:   "sync",
		Usage:  "Sync the .env file with .env.example",
		Action: commands.SyncEnv,
	}

	laravelClearCmd := &cli.Command{
		Name:    "clear",
		Aliases: []string{"c"},
		Usage:   "Clear all Laravel caches (cache, route, config, view)",
		Action:  commands.LaravelClear,
	}

	laravelFreshCmd := &cli.Command{
		Name:    "fresh",
		Aliases: []string{"f"},
		Usage:   "Run migrate:fresh with seeding",
		Action:  commands.LaravelFresh,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-seed",
				Usage: "Skip database seeding",
			},
		},
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "db",
				Usage:       "Database management",
				Subcommands: []*cli.Command{dbCreateCmd, dbListCmd, dbOpenCmd},
			},
			// Top-level commands with colon notation (reuse command definitions)
			{
				Name:    "db:create",
				Aliases: []string{"createdb"},
				Usage:   dbCreateCmd.Usage,
				Action:  dbCreateCmd.Action,
			},
			{
				Name:    "db:list",
				Aliases: []string{"listdb"},
				Usage:   dbListCmd.Usage,
				Action:  dbListCmd.Action,
			},
			{
				Name:    "db:open",
				Aliases: []string{"opendb"},
				Usage:   dbOpenCmd.Usage,
				Action:  dbOpenCmd.Action,
			},
			{
				Name:        "env",
				Usage:       "Environment management",
				Subcommands: []*cli.Command{envSqliteCmd, envMailtrapCmd, envMaildevCmd, envSyncCmd},
			},
			// Top-level commands with colon notation (reuse command definitions)
			{
				Name:   "env:sqlite",
				Usage:  envSqliteCmd.Usage,
				Action: envSqliteCmd.Action,
			},
			{
				Name:   "env:mailtrap",
				Usage:  envMailtrapCmd.Usage,
				Action: envMailtrapCmd.Action,
			},
			{
				Name:   "env:maildev",
				Usage:  envMaildevCmd.Usage,
				Action: envMaildevCmd.Action,
			},
			{
				Name:        "env:sync",
				Aliases:     []string{"sync:env"},
				Usage:       envSyncCmd.Usage,
				Description: "Sync the .env file with .env.example",
				Action:      envSyncCmd.Action,
			},
			{
				Name:    "config:edit",
				Aliases: []string{"edit:config"},
				Usage:   "Edit the Mortimer config file",
				Action:  commands.EditConfig,
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
			{
				Name:        "l",
				Usage:       "Laravel specific commands",
				Subcommands: []*cli.Command{laravelClearCmd, laravelFreshCmd},
			},
			// Top-level commands with colon notation (reuse command definitions)
			{
				Name:    "l:clear",
				Aliases: []string{"lc"},
				Usage:   laravelClearCmd.Usage,
				Action:  laravelClearCmd.Action,
			},
			{
				Name:    "l:fresh",
				Aliases: []string{"lf"},
				Usage:   laravelFreshCmd.Usage,
				Action:  laravelFreshCmd.Action,
				Flags:   laravelFreshCmd.Flags,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
