package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func ConnectSSH(cliContext *cli.Context) error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configFile := filepath.Join(homeDir, ".ssh", "config")

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("SSH config file not found: %s", configFile)
	}

	// Build the command: ssh $(grep "^Host " ~/.ssh/config | sed 's/Host //' | fzf)
	// We'll use a shell to execute this pipeline
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ssh $(grep '^Host ' %s | sed 's/Host //' | fzf)", configFile))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
