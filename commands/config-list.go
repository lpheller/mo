package commands

import (
	"fmt"
	"log"
	"mo/config"
	"sort"

	"github.com/urfave/cli/v2"
)

func ListConfigKeys(cliContext *cli.Context) error {
	// Load configuration file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return err
	}

	// Get all keys from ConfigPaths
	keys := make([]string, 0, len(cfg.ConfigPaths))
	for key := range cfg.ConfigPaths {
		keys = append(keys, key)
	}

	// Sort keys alphabetically
	sort.Strings(keys)

	if len(keys) == 0 {
		fmt.Println("No config keys found. Add keys to your config file using 'config_paths'.")
		return nil
	}

	fmt.Println("Available config keys:")
	fmt.Println("----------------------------")
	for _, key := range keys {
		fmt.Printf("  %s\n", key)
	}
	fmt.Println("----------------------------")
	fmt.Printf("Total: %d key(s)\n", len(keys))
	fmt.Println("\nUsage: mo qc <key>")

	return nil
}
