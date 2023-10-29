package utils

import (
	"fmt"
	"os"
	"regexp"
)

// ReplaceStringInFile replaces occurrences of oldStr with newStr in a file located at filePath
func ReplaceStringInFile(filePath, oldStr, newStr string) error {
	// Read the file into a byte slice
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Couldn't read file: %w", err)
	}

	// Convert the byte slice to a string
	fileContent := string(fileBytes)

	fmt.Println(fileContent)

	// Prepare the regular expression to match instances of oldStr not prefixed by a '#'
	re := regexp.MustCompile(`(?m)([^#])` + oldStr)

	// Replace the string using the regular expression
	updatedContent := re.ReplaceAllString(fileContent, "${1}"+newStr)

	// Write the updated content back to the file
	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("Couldn't write file: %w", err)
	}

	return nil
}

// func main() {
// 	err := ReplaceStringInFile("your-file.txt", "oldString", "newString")
// 	if err != nil {
// 		fmt.Println("An error occurred:", err)
// 	}
// }
