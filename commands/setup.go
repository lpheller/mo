package commands

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

// CheckProject will check the current directory for composer.json, artisan, and package.json files, then run appropriate commands.
func CheckProject(c *cli.Context) error {
	// Check if composer.json exists
	if _, err := os.Stat("composer.json"); err == nil {
		log.Println("composer.json found, running composer install...")
		if err := runCommand("composer", "install"); err != nil {
			return err
		}
	} else {
		log.Println("composer.json not found")
	}

	// Check if artisan exists (Laravel project)
	if _, err := os.Stat("artisan"); err == nil {
		log.Println("Laravel project detected...")

		// Check for .env file and copy .env.example if needed
		if err := ensureEnvFile(); err != nil {
			return err
		}

		// Check if APP_KEY is already set in .env
		hasAppKey, err := hasEnvAppKey()
		if err != nil {
			return err
		}

		// Run key:generate only if no APP_KEY is set
		if !hasAppKey {
			log.Println("No APP_KEY found, generating one...")
			if err := runCommand("php", "artisan", "key:generate"); err != nil {
				return err
			}
		} else {
			log.Println("APP_KEY already exists, skipping key:generate")
		}

		// Run artisan migrate and db:seed
		if err := runCommand("php", "artisan", "migrate"); err != nil {
			return err
		}
		if err := runCommand("php", "artisan", "db:seed"); err != nil {
			return err
		}
	} else {
		log.Println("artisan file not found")
	}

	// Check if package.json exists
	if _, err := os.Stat("package.json"); err == nil {
		log.Println("package.json found, running npm install...")
		if err := runCommand("npm", "install"); err != nil {
			return err
		}

		// Check if build script exists in package.json
		hasBuildScript, err := hasNpmScript("build")
		if err != nil {
			return err
		}

		if hasBuildScript {
			log.Println("Build script found, running npm run build...")
			if err := runCommand("npm", "run", "build"); err != nil {
				return err
			}
		} else {
			log.Println("No build script found")
		}
	} else {
		log.Println("package.json not found")
	}

	return nil
}

// ensureEnvFile checks if .env exists, and if not, copies .env.example to .env
func ensureEnvFile() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println(".env file not found, copying from .env.example...")

		if _, err := os.Stat(".env.example"); err == nil {
			srcFile, err := os.Open(".env.example")
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(".env")
			if err != nil {
				return err
			}
			defer destFile.Close()

			if _, err := io.Copy(destFile, srcFile); err != nil {
				return err
			}

			log.Println(".env file created from .env.example")
		} else {
			log.Println(".env.example not found, cannot create .env")
			return err
		}
	} else {
		log.Println(".env file exists")
	}

	return nil
}

// hasEnvAppKey checks if the APP_KEY is set in the .env file
func hasEnvAppKey() (bool, error) {
	file, err := os.Open(".env")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "APP_KEY=") && len(line) > len("APP_KEY=") {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// runCommand is a helper function to run system commands
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}

// hasNpmScript checks if a given npm script exists in package.json
func hasNpmScript(script string) (bool, error) {
	file, err := os.Open("package.json")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), `"scripts": {`) {
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), fmt.Sprintf(`"%s":`, script)) {
					return true, nil
				}
				// Stop when reaching the end of the scripts section
				if strings.Contains(scanner.Text(), `},`) {
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
