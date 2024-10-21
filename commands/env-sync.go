package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

/**
 * Sync between the .env.example and .env file
 *
 * @param c *cli.Context
 * @return error
 */
func SyncEnv(c *cli.Context) error {
	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		file, err := os.Create(".env.example")
		if err != nil {
			return fmt.Errorf("error creating .env.example: %v", err)
		}
		defer file.Close()
		fmt.Println(".env.example file created.")
	}

	envVariables, err := getEnvVariablesFrom(".env")
	if err != nil {
		return fmt.Errorf("error reading .env: %v", err)
	}

	exampleVariables, err := getEnvVariablesFrom(".env.example")
	if err != nil {
		return fmt.Errorf("error reading .env.example: %v", err)
	}

	var difference []string
	for _, variable := range envVariables {
		if !contains(exampleVariables, variable) {
			difference = append(difference, variable)
		}
	}

	file, err := os.OpenFile(".env.example", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening .env.example: %v", err)
	}
	defer file.Close()

	if len(difference) == 0 {
		fmt.Println("No missing variables found")
		return nil
	}

	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("error writing newline to .env.example: %v", err)
	}

	for _, variable := range difference {
		line := fmt.Sprintf("%s=\n", variable)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("error writing variable to .env.example: %v", err)
		}
	}

	fmt.Println("Added missing variables to .env.example:")
	for _, variable := range difference {
		fmt.Println(variable)
	}

	return nil
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func getEnvVariablesFrom(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var envVariables []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) >= 1 {
			envVariables = append(envVariables, strings.TrimSpace(parts[0]))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envVariables, nil
}
