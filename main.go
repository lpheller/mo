package main

import (
	// "fmt"

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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
