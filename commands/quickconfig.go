package commands

import (
	"fmt"
	"log"
	"mo/config" // Adjust import path if necessary
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func QuickConfig(cliContext *cli.Context) error {
	// Load configuration file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return err
	}

	// Ensure a config key argument is provided
	if cliContext.Args().Len() == 0 {
		return fmt.Errorf("missing config key: you must provide a config key (e.g., mortimer, nvim, git)")
	}

	key := cliContext.Args().Get(0)
	path, exists := cfg.ConfigPaths[key]
	if !exists {
		fmt.Printf("No config settings for '%s'\n", key)
		return nil
	}

	// Use editor from config or override with --editor flag
	editor := cfg.Editor
	if cliContext.IsSet("editor") {
		editor = cliContext.String("editor")
		fmt.Printf("Using editor from flag: %s\n", editor)
	} else {
		fmt.Printf("Using editor from config: %s\n", editor)
	}

	// Open the config file using the specified editor
	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error opening config file: %v", err)
		return err
	}

	fmt.Printf("Opened %s using %s\n", path, editor)
	return nil
}
