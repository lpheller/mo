package commands

import (
	"fmt"
	"mo/utils"

	"github.com/urfave/cli/v2"
)

func LaravelClear(cliContext *cli.Context) error {
	if !fileExists("artisan") {
		return fmt.Errorf("not a Laravel project")
	}

	commands := [][]string{
		{"php", "artisan", "cache:clear"},
		{"php", "artisan", "route:clear"},
		{"php", "artisan", "config:clear"},
		{"php", "artisan", "view:clear"},
	}

	for _, cmd := range commands {
		if err := utils.RunCommand(cmd[0], cmd[1:]...); err != nil {
			return fmt.Errorf("error running %s: %w", cmd[2], err)
		}
	}

	fmt.Println("✓ All Laravel caches cleared!")
	return nil
}
