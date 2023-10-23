package main

import (
	// "fmt"
	"log"
	"os"

	"mogo/commands"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "db:open",
				Aliases: []string{"opendb"},
				Usage:   "Open the database in the default editor",
				Action:  commands.OpenDatabase,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
