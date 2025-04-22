package commands

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mo/utils"

	"github.com/urfave/cli/v2"
)

func Push(cliContext *cli.Context) error {
	localEnv, err := utils.EnsureRequiredEnvVars("push")
	if err != nil {
		return err
	}

	remoteEnvPath := fmt.Sprintf("%s/.env", localEnv["PUSH_PROJECT_DIR"])
	log.Printf("Remote environment path: %s", remoteEnvPath)
	remoteAppEnv, err := getRemoteEnvValue(localEnv, remoteEnvPath, "APP_ENV", "push")
	if err != nil {
		return fmt.Errorf("error fetching remote APP_ENV: %v", err)
	}

	if remoteAppEnv != "staging" {
		fmt.Println("Remote environment is not set to 'staging'. Aborting to prevent overwriting data.")
		return nil
	}

	if !cliContext.Bool("storage") && !cliContext.Bool("database") {
		fmt.Println("No flags set. Use --storage or --database to push data.")
		return nil
	}

	if cliContext.Bool("storage") {
		if err := pushStorage(localEnv); err != nil {
			return err
		}
	}

	if cliContext.Bool("database") {
		if err := pushDatabase(localEnv); err != nil {
			return err
		}
	}

	return nil
}

func pushStorage(env map[string]string) error {
	fmt.Println("Pushing storage folder...")

	dateOutput, err := exec.Command("date", "+%F").Output()
	if err != nil {
		return fmt.Errorf("error getting date: %v", err)
	}
	storageFile := fmt.Sprintf("storage_%s.tar.gz", strings.TrimSpace(string(dateOutput)))

	// Compress the local storage folder using native Go
	fmt.Println("Compressing local storage folder...")
	if err := compressFolder("storage/app/public", storageFile); err != nil {
		return fmt.Errorf("error compressing local storage folder: %v", err)
	}

	localPath := storageFile
	remotePath := fmt.Sprintf("%s@%s:/tmp/%s", env["PUSH_SSH_USER"], env["PUSH_HOST"], storageFile)

	if err := utils.RunCommand("scp", localPath, remotePath); err != nil {
		return fmt.Errorf("error uploading storage file: %v", err)
	}

	remoteCmd := fmt.Sprintf("cd %s/storage/app/public && tar -xzf /tmp/%s && rm /tmp/%s", env["PUSH_PROJECT_DIR"], storageFile, storageFile)
	if err := utils.RunRemoteCommand(env["PUSH_SSH_USER"], env["PUSH_HOST"], remoteCmd); err != nil {
		return fmt.Errorf("error extracting storage file on remote: %v", err)
	}

	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("error deleting local storage file: %v", err)
	}

	fmt.Println("Storage folder successfully pushed!")
	return nil
}

func compressFolder(sourceDir, outputFile string) error {
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = filepath.Walk(sourceDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through file: %v", err)
		}

		if fi.Name() == ".DS_Store" {
			return nil
		}

		relPath := strings.TrimPrefix(file, sourceDir)
		relPath = strings.TrimPrefix(relPath, string(filepath.Separator))

		if relPath == "" {
			relPath = "."
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return fmt.Errorf("error creating tar header: %v", err)
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("error writing tar header: %v", err)
		}

		if !fi.IsDir() {
			fileContent, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("error opening file: %v", err)
			}
			defer fileContent.Close()

			if _, err := io.Copy(tarWriter, fileContent); err != nil {
				return fmt.Errorf("error writing file content to tar archive: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error compressing folder: %v", err)
	}

	return nil
}

func pushDatabase(env map[string]string) error {
	fmt.Println("Pushing database...")
	envManager := utils.NewEnvManager(".env")
	localDBName, _, _ := envManager.GetVar("DB_DATABASE")
	localDBUser, _, _ := envManager.GetVar("DB_USERNAME")
	localDBPassword, _, _ := envManager.GetVar("DB_PASSWORD")

	log.Printf("Local DB Name: %s", localDBName)
	log.Printf("Local DB User: %s", localDBUser)
	log.Printf("Local DB Password: %s", localDBPassword)

	if localDBName == "" || localDBUser == "" {
		return fmt.Errorf("local database credentials are missing in .env file")
	}

	dumpFile := fmt.Sprintf("%s-dump.sql", localDBName)

	var dumpCmd string
	if localDBPassword != "" {
		dumpCmd = fmt.Sprintf("mysqldump -u %s -p%s %s > %s", localDBUser, localDBPassword, localDBName, dumpFile)
	} else {
		dumpCmd = fmt.Sprintf("mysqldump -u %s %s > %s", localDBUser, localDBName, dumpFile)
	}
	if err := utils.RunCommand("bash", "-c", dumpCmd); err != nil {
		return fmt.Errorf("error creating local database dump: %v", err)
	}

	remotePath := fmt.Sprintf("%s@%s:/tmp/%s", env["PUSH_SSH_USER"], env["PUSH_HOST"], dumpFile)

	if err := utils.RunCommand("scp", dumpFile, remotePath); err != nil {
		return fmt.Errorf("error uploading database dump: %v", err)
	}

	remoteDBName, err := getRemoteEnvValue(env, fmt.Sprintf("%s/.env", env["PUSH_PROJECT_DIR"]), "DB_DATABASE", "push")
	if err != nil {
		return fmt.Errorf("error fetching remote DB name: %v", err)
	}
	remoteDBUser, err := getRemoteEnvValue(env, fmt.Sprintf("%s/.env", env["PUSH_PROJECT_DIR"]), "DB_USERNAME", "push")
	if err != nil {
		return fmt.Errorf("error fetching remote DB user: %v", err)
	}
	remoteDBPassword, err := getRemoteEnvValue(env, fmt.Sprintf("%s/.env", env["PUSH_PROJECT_DIR"]), "DB_PASSWORD", "push")
	if err != nil {
		return fmt.Errorf("error fetching remote DB password: %v", err)
	}

	remoteCmd := fmt.Sprintf("mysql -u %s -p%s %s < /tmp/%s && rm /tmp/%s", remoteDBUser, remoteDBPassword, remoteDBName, dumpFile, dumpFile)
	if err := utils.RunRemoteCommand(env["PUSH_SSH_USER"], env["PUSH_HOST"], remoteCmd); err != nil {
		return fmt.Errorf("error importing database dump on remote: %v", err)
	}

	if err := os.Remove(dumpFile); err != nil {
		return fmt.Errorf("error deleting local database dump: %v", err)
	}

	fmt.Println("Database successfully pushed!")
	return nil
}
