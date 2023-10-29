package commands

import (
	"bufio"
	"fmt"
	"mo/utils"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

func checkFileForString(filePath string, searchString string) (bool, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, -1, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.Contains(line, searchString) {
			return true, lineNumber, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, -1, err
	}

	return false, -1, nil
}

func EnvSqlite(c *cli.Context) error {
	filePath := ".env"

	utils.ReplaceStringInFile(filePath, "DB_", "#DB_")
	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return fmt.Errorf("Couldn't read file: %w", err)
	}
	fileContent := string(fileBytes)

	// This will match "#DB_CONNECTION=" followed by any characters to the end of the line
	re := regexp.MustCompile(`#DB_CONNECTION=.*`)
	newConfig := re.ReplaceAllString(fileContent, "DB_CONNECTION=sqlite")

	err = os.WriteFile(filePath, []byte(newConfig), 0644)
	return nil

}
