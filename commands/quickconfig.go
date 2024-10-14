package commands

import (
	"fmt"
	"log"
	"mo/config" // Adjust import path if necessary
	"os/exec"

	"github.com/urfave/cli/v2"
)

func QuickConfig(c *cli.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return err
	}

	// Ensure a key is provided
	if c.Args().Len() == 0 {
		log.Fatalf("You must provide a config key (e.g., mortimer, nvim, git)")
		return fmt.Errorf("missing config key")
	}

	key := c.Args().Get(0)
	path, exists := cfg.ConfigPaths[key]
	if !exists {
		fmt.Printf("No config settings for '%s'\n", key)
		return nil
	}

	// Choose the editor from the config (default to nvim if not specified)
	editor := cfg.Editor
	if len(editor) == 0 {
		editor = "nvim"
	}

	// Open the file with the specified editor
	cmd := exec.Command(editor, path)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return err
	}

	fmt.Printf("Opened %s using %s\n", path, editor)
	return nil
}
