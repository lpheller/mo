package commands

import (
	"fmt"
	"mo/utils"

	"github.com/urfave/cli/v2"
)

func LaravelFresh(cliContext *cli.Context) error {
	if !fileExists("artisan") {
		return fmt.Errorf("not a Laravel project")
	}

	args := []string{"artisan", "migrate:fresh"}

	if !cliContext.Bool("no-seed") {
		args = append(args, "--seed")
		fmt.Println("Running migrate:fresh --seed...")
	} else {
		fmt.Println("Running migrate:fresh...")
	}

	if err := utils.RunCommand("php", args...); err != nil {
		return fmt.Errorf("error running migrate:fresh: %w", err)
	}

	fmt.Println("✓ Database refreshed successfully!")
	return nil
}
