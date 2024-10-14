package commands

import (
	"fmt"
	"log"
	"mo/config" // Adjust import path if necessary
	"os"
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
	// Dump the --editor flag value
	fmt.Printf("Flag --editor value: '%s'\n", c.String("editor"))

	// If the --editor flag is provided, use it instead
	if c.IsSet("editor") {
		// Debugging: Print the flag value
		fmt.Printf("Using editor from flag: %s\n", c.String("editor"))
		editor = c.String("editor")
	} else {
		// Debugging: Print which editor is being used
		fmt.Printf("Using editor from config: %s\n", editor)
	}

	// output the editor
	// fmt.Printf("Opening %s using %s\n", key, editor)

	// Open the file with the specified editor
	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return err
	}

	fmt.Printf("Opened %s using %s\n", path, editor)
	return nil
}
