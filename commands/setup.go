package commands

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"mo/utils"

	"github.com/urfave/cli/v2"
)

func CheckProject(cliContext *cli.Context) error {
	if err := handleComposer(); err != nil {
		return err
	}

	if err := handleLaravel(); err != nil {
		return err
	}

	if err := handleNode(); err != nil {
		return err
	}

	return nil
}

func handleComposer() error {
	if fileExists("composer.json") {
		log.Println("composer.json found, running composer install...")
		if err := utils.RunCommand("composer", "install"); err != nil {
			return fmt.Errorf("composer install failed: %w", err)
		}
	} else {
		log.Println("composer.json not found")
	}
	return nil
}

func handleLaravel() error {
	if !fileExists("artisan") {
		log.Println("artisan file not found")
		return nil
	}

	log.Println("Laravel project detected...")

	if err := ensureEnvFile(); err != nil {
		return fmt.Errorf("error ensuring .env file: %w", err)
	}

	hasAppKey, err := hasEnvAppKey()
	if err != nil {
		return fmt.Errorf("error checking APP_KEY: %w", err)
	}
	if !hasAppKey {
		log.Println("No APP_KEY found, generating one...")
		if err := utils.RunCommand("php", "artisan", "key:generate"); err != nil {
			return fmt.Errorf("key:generate failed: %w", err)
		}
	} else {
		log.Println("APP_KEY already exists, skipping key:generate")
	}

	if err := utils.RunCommand("php", "artisan", "migrate"); err != nil {
		return fmt.Errorf("artisan migrate failed: %w", err)
	}
	if err := utils.RunCommand("php", "artisan", "db:seed"); err != nil {
		return fmt.Errorf("artisan db:seed failed: %w", err)
	}

	return nil
}

func handleNode() error {
	if !fileExists("package.json") {
		log.Println("package.json not found")
		return nil
	}

	log.Println("package.json found, running npm install...")
	if err := utils.RunCommand("npm", "install"); err != nil {
		return fmt.Errorf("npm install failed: %w", err)
	}

	hasBuildScript, err := hasNpmScript("build")
	if err != nil {
		return fmt.Errorf("error checking build script: %w", err)
	}
	if hasBuildScript {
		log.Println("Build script found, running npm run build...")
		if err := utils.RunCommand("npm", "run", "build"); err != nil {
			return fmt.Errorf("npm run build failed: %w", err)
		}
	} else {
		log.Println("No build script found")
	}

	return nil
}

func ensureEnvFile() error {
	if fileExists(".env") {
		log.Println(".env file exists")
		return nil
	}

	log.Println(".env file not found, copying from .env.example...")
	if !fileExists(".env.example") {
		return fmt.Errorf(".env.example not found")
	}

	if err := copyFile(".env.example", ".env"); err != nil {
		return fmt.Errorf("error copying .env.example to .env: %w", err)
	}

	log.Println(".env file created from .env.example")
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

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

func hasNpmScript(script string) (bool, error) {
	file, err := os.Open("package.json")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `"scripts": {`) {
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), fmt.Sprintf(`"%s":`, script)) {
					return true, nil
				}
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
