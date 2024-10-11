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
	// Check if .env.example exists; if not, create it
	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		file, err := os.Create(".env.example")
		if err != nil {
			fmt.Println("Error creating .env.example:", err)
			return nil
		}
		defer file.Close()
		fmt.Println(".env.example file created.")
	}

	envVariables, err := getEnvVariablesFrom(".env")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	exampleVariables, err := getEnvVariablesFrom(".env.example")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	difference := []string{}
	for _, variable := range envVariables {
		if !contains(exampleVariables, variable) {
			difference = append(difference, variable)
		}
	}

	file, err := os.OpenFile(".env.example", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	if len(difference) == 0 {
		fmt.Println("No missing Variables found")
		return nil
	}

	line := fmt.Sprintf("\n") // Add a newline before the variable
	_, err = file.WriteString(line)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, variable := range difference {
		line := fmt.Sprintf("%s=\n", variable)
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	fmt.Println("Added missing Variables to .env.example:")
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
