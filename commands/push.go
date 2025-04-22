package commands

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mo/utils"

	"github.com/urfave/cli/v2"
)

func Push(cliContext *cli.Context) error {
	localEnv, err := loadLocalEnv()
	if err != nil {
		return err
	}

	// Check if the remote environment is staging
	remoteEnvPath := fmt.Sprintf("%s/.env", localEnv["REMOTE_PROJECT_DIR"])
	remoteAppEnv, err := getRemoteEnvValue(localEnv, remoteEnvPath, "APP_ENV")
	if err != nil {
		return fmt.Errorf("error fetching remote APP_ENV: %v", err)
	}

	if remoteAppEnv != "staging" {
		fmt.Println("Remote environment is not set to 'staging'. Aborting to prevent overwriting data.")
		return nil
	}

	// If no flags are set, print usage and return
	if !cliContext.Bool("storage") && !cliContext.Bool("database") {
		fmt.Println("No flags set. Use --storage or --database to push data.")
		return nil
	}

	// Check if the --storage flag is set, and push the storage folder if true.
	if cliContext.Bool("storage") {
		if err := pushStorage(localEnv); err != nil {
			return err
		}
	}

	// Check if the --database flag is set, and push the database if true.
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
	remotePath := fmt.Sprintf("%s@%s:/tmp/%s", env["REMOTE_SSH_USER"], env["REMOTE_HOST"], storageFile)

	// Upload the compressed file to the remote server
	if err := utils.RunCommand("scp", localPath, remotePath); err != nil {
		return fmt.Errorf("error uploading storage file: %v", err)
	}

	remoteCmd := fmt.Sprintf("cd %s/storage/app/public && tar -xzf /tmp/%s && rm /tmp/%s", env["REMOTE_PROJECT_DIR"], storageFile, storageFile)
	if err := utils.RunRemoteCommand(env["REMOTE_SSH_USER"], env["REMOTE_HOST"], remoteCmd); err != nil {
		return fmt.Errorf("error extracting storage file on remote: %v", err)
	}

	// Remove the local compressed file
	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("error deleting local storage file: %v", err)
	}

	fmt.Println("Storage folder successfully pushed!")
	return nil
}

// compressFolder compresses a folder into a .tar.gz file using native Go libraries
func compressFolder(sourceDir, outputFile string) error {
	// Create the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Walk through the source directory and add files to the tar archive
	err = filepath.Walk(sourceDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through file: %v", err)
		}

		// Skip unwanted files like .DS_Store
		if fi.Name() == ".DS_Store" {
			return nil
		}

		// Get the relative path to maintain folder structure
		relPath := strings.TrimPrefix(file, sourceDir)
		relPath = strings.TrimPrefix(relPath, string(filepath.Separator))

		// Ensure the relative path is not empty
		if relPath == "" {
			relPath = "."
		}

		// Create a tar header for the file
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return fmt.Errorf("error creating tar header: %v", err)
		}
		header.Name = relPath

		// Write the header to the tar archive
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("error writing tar header: %v", err)
		}

		// If the file is not a directory, write its content to the tar archive
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

	localEnv, err := loadLocalEnv()
	if err != nil {
		return err
	}

	localDBName := localEnv["DB_DATABASE"]
	localDBUser := localEnv["DB_USERNAME"]
	localDBPassword := localEnv["DB_PASSWORD"]

	if localDBName == "" || localDBUser == "" || localDBPassword == "" {
		return fmt.Errorf("local database credentials are missing in .env file")
	}

	dumpFile := fmt.Sprintf("%s-dump.sql", localDBName)

	// Create a database dump locally
	dumpCmd := fmt.Sprintf("mysqldump -u %s -p%s %s > %s", localDBUser, localDBPassword, localDBName, dumpFile)
	if err := utils.RunCommand("bash", "-c", dumpCmd); err != nil {
		return fmt.Errorf("error creating local database dump: %v", err)
	}

	remotePath := fmt.Sprintf("%s@%s:/tmp/%s", env["REMOTE_SSH_USER"], env["REMOTE_HOST"], dumpFile)

	// Upload the dump file to the remote server
	if err := utils.RunCommand("scp", dumpFile, remotePath); err != nil {
		return fmt.Errorf("error uploading database dump: %v", err)
	}

	remoteDBName, err := getRemoteEnvValue(env, fmt.Sprintf("%s/.env", env["REMOTE_PROJECT_DIR"]), "DB_DATABASE")
	if err != nil {
		return fmt.Errorf("error fetching remote DB name: %v", err)
	}
	remoteDBUser, err := getRemoteEnvValue(env, fmt.Sprintf("%s/.env", env["REMOTE_PROJECT_DIR"]), "DB_USERNAME")
	if err != nil {
		return fmt.Errorf("error fetching remote DB user: %v", err)
	}
	remoteDBPassword, err := getRemoteEnvValue(env, fmt.Sprintf("%s/.env", env["REMOTE_PROJECT_DIR"]), "DB_PASSWORD")
	if err != nil {
		return fmt.Errorf("error fetching remote DB password: %v", err)
	}

	// Import the dump file into the remote database
	remoteCmd := fmt.Sprintf("mysql -u %s -p%s %s < /tmp/%s && rm /tmp/%s", remoteDBUser, remoteDBPassword, remoteDBName, dumpFile, dumpFile)
	if err := utils.RunRemoteCommand(env["REMOTE_SSH_USER"], env["REMOTE_HOST"], remoteCmd); err != nil {
		return fmt.Errorf("error importing database dump on remote: %v", err)
	}

	// Remove the local dump file
	if err := os.Remove(dumpFile); err != nil {
		return fmt.Errorf("error deleting local database dump: %v", err)
	}

	fmt.Println("Database successfully pushed!")
	return nil
}
