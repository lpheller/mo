package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// getSSHConfigPath returns the path to the SSH config file
func getSSHConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".ssh", "config"), nil
}

// parseConnectionString extracts user and host from a connection string (user@host or user@host:port)
func parseConnectionString(connection string) (user, host string, err error) {
	// Handle port in connection string (user@host:port)
	parts := strings.Split(connection, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid connection format: expected 'user@host', got '%s'", connection)
	}

	user = parts[0]
	hostWithPort := parts[1]

	// Remove port if present (host:port)
	hostParts := strings.Split(hostWithPort, ":")
	host = hostParts[0]

	return user, host, nil
}

// addSSHEntryInternal adds an SSH entry to the config file
func addSSHEntryInternal(name, user, host string) error {
	configFile, err := getSSHConfigPath()
	if err != nil {
		return err
	}

	// Check if the entry already exists
	if err := checkSSHEntryExists(configFile, name); err != nil {
		return err
	}

	// Get home directory for .ssh directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Ensure .ssh directory exists
	sshDir := filepath.Join(homeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	// Check if file exists to fix permissions if needed
	fileExists := true
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fileExists = false
	}

	// Open file for appending (create if it doesn't exist)
	file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open SSH config file: %w", err)
	}
	defer file.Close()

	// Fix permissions if file already existed (OpenFile only sets permissions on creation)
	if fileExists {
		if err := os.Chmod(configFile, 0600); err != nil {
			return fmt.Errorf("failed to set correct permissions on SSH config file: %w", err)
		}
	}

	// Append the new entry
	entry := fmt.Sprintf("\nHost %s\n  HostName %s\n  User %s\n", name, host, user)
	if _, err := file.WriteString(entry); err != nil {
		return fmt.Errorf("failed to write SSH entry: %w", err)
	}

	return nil
}

func AddSSHEntry(cliContext *cli.Context) error {
	if cliContext.Args().Len() != 2 {
		return fmt.Errorf("usage: ssh:add <name> <connection>")
	}

	name := cliContext.Args().Get(0)
	connection := cliContext.Args().Get(1)

	// Extract user and host from connection string
	user, host, err := parseConnectionString(connection)
	if err != nil {
		return err
	}

	// Add the entry
	if err := addSSHEntryInternal(name, user, host); err != nil {
		return err
	}

	fmt.Printf("Added SSH entry for '%s' with connection '%s'.\n", name, connection)
	return nil
}

func checkSSHEntryExists(configFile, name string) error {
	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// File doesn't exist, so entry doesn't exist either
		return nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		// If we can't read it, we'll try to create it anyway
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Check for exact match: "Host <name>" (with optional whitespace)
		if strings.HasPrefix(line, "Host ") {
			hostName := strings.TrimSpace(strings.TrimPrefix(line, "Host"))
			if hostName == name {
				return fmt.Errorf("entry for '%s' already exists in %s", name, configFile)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading SSH config file: %w", err)
	}

	return nil
}

// SSHHostEntry represents an SSH host entry from the config
type SSHHostEntry struct {
	Name     string
	HostName string
	User     string
}

// parseSSHConfig parses the SSH config file and returns all host entries
func parseSSHConfig(configFile string) ([]SSHHostEntry, error) {
	var entries []SSHHostEntry

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return entries, nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return entries, fmt.Errorf("failed to open SSH config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentEntry *SSHHostEntry

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for Host declaration
		if strings.HasPrefix(line, "Host ") {
			// Save previous entry if exists
			if currentEntry != nil {
				entries = append(entries, *currentEntry)
			}

			hostName := strings.TrimSpace(strings.TrimPrefix(line, "Host"))
			// Skip wildcard entries
			if strings.Contains(hostName, "*") || strings.Contains(hostName, "?") {
				currentEntry = nil
				continue
			}

			currentEntry = &SSHHostEntry{
				Name: hostName,
			}
		} else if currentEntry != nil {
			// Parse HostName and User
			if strings.HasPrefix(line, "HostName ") {
				currentEntry.HostName = strings.TrimSpace(strings.TrimPrefix(line, "HostName"))
			} else if strings.HasPrefix(line, "User ") {
				currentEntry.User = strings.TrimSpace(strings.TrimPrefix(line, "User"))
			}
		}
	}

	// Don't forget the last entry
	if currentEntry != nil {
		entries = append(entries, *currentEntry)
	}

	if err := scanner.Err(); err != nil {
		return entries, fmt.Errorf("error reading SSH config file: %w", err)
	}

	return entries, nil
}

func findSSHEntryByConnection(connection string) (string, error) {
	user, host, err := parseConnectionString(connection)
	if err != nil {
		return "", err
	}

	configFile, err := getSSHConfigPath()
	if err != nil {
		return "", err
	}

	entries, err := parseSSHConfig(configFile)
	if err != nil {
		return "", err
	}

	// Search for matching entry
	for _, entry := range entries {
		if entry.HostName == host && entry.User == user {
			return entry.Name, nil
		}
	}

	return "", nil // Not found, but no error
}
