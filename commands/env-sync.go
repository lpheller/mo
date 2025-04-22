package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mo/utils"

	"github.com/urfave/cli/v2"
)

func SyncEnv(cliContext *cli.Context) error {
	// envManager := utils.NewEnvManager(".env")
	exampleManager := utils.NewEnvManager(".env.example")

	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		if err := os.WriteFile(".env.example", []byte{}, 0644); err != nil {
			return fmt.Errorf("error creating .env.example: %v", err)
		}
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

	for i := 0; i < len(difference); i++ {
		if strings.TrimSpace(difference[i]) == "" {
			difference = append(difference[:i], difference[i+1:]...)
			i--
		}
	}

	differenceMap := make(map[string]bool)
	for _, variable := range difference {
		differenceMap[variable] = true
	}
	difference = []string{}
	for variable := range differenceMap {
		difference = append(difference, variable)
	}

	for _, variable := range difference {
		fmt.Printf("Adding %s to .env.example\n", variable)
		if err := exampleManager.SetVar(variable, ""); err != nil {
			return fmt.Errorf("error adding variable %s to .env.example: %v", variable, err)
		}
	}

	if len(difference) == 0 {
		fmt.Println("No missing variables found")
	} else {
		fmt.Println("Added missing variables to .env.example:")
		for _, variable := range difference {
			fmt.Println(variable)
		}
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
