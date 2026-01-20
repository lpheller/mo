package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func ListSSHEntries(cliContext *cli.Context) error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configFile := filepath.Join(homeDir, ".ssh", "config")

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("No SSH entries found in %s.\n", configFile)
		return nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("failed to open SSH config file: %w", err)
	}
	defer file.Close()

	var entries []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Host ") {
			hostName := strings.TrimSpace(strings.TrimPrefix(line, "Host"))
			// Skip wildcard entries and patterns
			if !strings.Contains(hostName, "*") && !strings.Contains(hostName, "?") {
				entries = append(entries, hostName)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading SSH config file: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("No SSH entries found in %s.\n", configFile)
		return nil
	}

	fmt.Println("SSH Entries:")
	fmt.Println("----------------------------")
	for _, entry := range entries {
		fmt.Printf("  %s\n", entry)
	}
	fmt.Println("----------------------------")
	fmt.Printf("Total: %d entry(ies)\n", len(entries))

	return nil
}
