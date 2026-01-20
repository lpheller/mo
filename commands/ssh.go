package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

// promptForHostName prompts the user for a host name and validates it
func promptForHostName() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a name for this SSH connection: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	name := strings.TrimSpace(input)
	if name == "" {
		return "", fmt.Errorf("host name cannot be empty")
	}

	// Validate name (no spaces, no special characters that would break SSH config)
	if strings.ContainsAny(name, " \t\n\r") {
		return "", fmt.Errorf("host name cannot contain spaces")
	}

	return name, nil
}

// connectOrAddSSHH connects to an SSH host, adding it to config if it doesn't exist
func connectOrAddSSHH(connection string) error {
	// Check if entry already exists
	existingName, err := findSSHEntryByConnection(connection)
	if err != nil {
		return fmt.Errorf("error checking for existing SSH entry: %w", err)
	}

	var hostName string

	if existingName != "" {
		// Entry exists, use it
		hostName = existingName
		fmt.Printf("Found existing SSH entry: %s\n", hostName)
	} else {
		// Entry doesn't exist, prompt for name and add it
		name, err := promptForHostName()
		if err != nil {
			return err
		}

		// Parse connection string
		user, host, err := parseConnectionString(connection)
		if err != nil {
			return err
		}

		// Add the entry
		if err := addSSHEntryInternal(name, user, host); err != nil {
			return err
		}

		fmt.Printf("Added SSH entry for '%s' with connection '%s'.\n", name, connection)
		hostName = name
	}

	// Connect via SSH
	fmt.Printf("Connecting to %s...\n", hostName)
	cmd := exec.Command("ssh", "-t", hostName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SSHCommand(cliContext *cli.Context) error {
	// If connection string is provided as argument, use connectOrAddSSHH
	if cliContext.Args().Len() > 0 {
		connection := cliContext.Args().Get(0)
		return connectOrAddSSHH(connection)
	}

	// Otherwise, use fzf for interactive selection
	return ConnectSSH(cliContext)
}
