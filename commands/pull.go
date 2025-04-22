package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"mo/config"
	"mo/utils"

	"github.com/urfave/cli/v2"
)

func Pull(cliContext *cli.Context) error {
	localEnv, err := utils.EnsureRequiredEnvVars("pull")
	if err != nil {
		return err
	}

	if !cliContext.Bool("storage") && !cliContext.Bool("database") {
		fmt.Println("No flags set. Use --storage or --database to pull data.")
		return nil
	}

	if cliContext.Bool("storage") {
		if err := pullStorage(localEnv); err != nil {
			return err
		}
	}

	if cliContext.Bool("database") {
		if err := pullDatabase(localEnv); err != nil {
			return err
		}
	}

	return nil
}

func pullStorage(env map[string]string) error {
	fmt.Println("Pulling storage folder...")

	dateOutput, err := exec.Command("date", "+%F").Output()
	if err != nil {
		return fmt.Errorf("error getting date: %v", err)
	}
	storageFile := fmt.Sprintf("storage_%s.tar.gz", strings.TrimSpace(string(dateOutput)))
	remoteCmd := fmt.Sprintf("cd %s && tar -cvzf /tmp/%s -C storage/app/public .", env["PULL_PROJECT_DIR"], storageFile)

	if err := utils.RunRemoteCommand(env["PULL_SSH_USER"], env["PULL_HOST"], remoteCmd); err != nil {
		return fmt.Errorf("error compressing storage folder on remote: %v", err)
	}

	localPath := fmt.Sprintf("/tmp/%s", storageFile)
	remotePath := fmt.Sprintf("%s@%s:/tmp/%s", env["PULL_SSH_USER"], env["PULL_HOST"], storageFile)

	if err := utils.RunCommand("scp", remotePath, localPath); err != nil {
		return fmt.Errorf("error downloading storage file: %v", err)
	}

	if err := utils.RunCommand("tar", "-xzf", localPath, "-C", "storage/app/public"); err != nil {
		return fmt.Errorf("error extracting storage file locally: %v", err)
	}

	if err := utils.RunRemoteCommand(env["PULL_SSH_USER"], env["PULL_HOST"], fmt.Sprintf("rm /tmp/%s", storageFile)); err != nil {
		return fmt.Errorf("error deleting remote storage file: %v", err)
	}

	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("error deleting local storage file: %v", err)
	}

	fmt.Println("Storage folder successfully pulled!")
	return nil
}

func pullDatabase(env map[string]string) error {
	fmt.Println("Pulling database...")

	remoteEnvPath := fmt.Sprintf("%s/.env", env["PULL_PROJECT_DIR"])
	fmt.Println("Remote env path:", remoteEnvPath)
	remoteDBName, err := getRemoteEnvValue(env, remoteEnvPath, "DB_DATABASE", "pull")
	if err != nil {
		return fmt.Errorf("error fetching remote DB name: %v", err)
	}
	fmt.Println("Remote DB Name: %s\n", remoteDBName)

	remoteDBUser, err := getRemoteEnvValue(env, remoteEnvPath, "DB_USERNAME", "pull")
	if err != nil {
		return err
	}
	fmt.Println("Remote DB User: %s\n", remoteDBUser)
	remoteDBPassword, err := getRemoteEnvValue(env, remoteEnvPath, "DB_PASSWORD", "pull")
	if err != nil {
		return err
	}

	dumpFile := fmt.Sprintf("/tmp/%s-dump.sql", remoteDBName)
	remoteCmd := fmt.Sprintf("mysqldump -u %s -p%s %s > %s", remoteDBUser, remoteDBPassword, remoteDBName, dumpFile)

	fmt.Println("Remote command to create dump:", remoteCmd)

	if err := utils.RunRemoteCommand(env["PULL_SSH_USER"], env["PULL_HOST"], remoteCmd); err != nil {
		return fmt.Errorf("error creating database dump on remote: %v", err)
	}

	localPath := fmt.Sprintf("%s-dump.sql", remoteDBName)
	remotePath := fmt.Sprintf("%s@%s:%s", env["PULL_SSH_USER"], env["PULL_HOST"], dumpFile)

	if err := utils.RunCommand("scp", remotePath, localPath); err != nil {
		return fmt.Errorf("error downloading database dump: %v", err)
	}

	localEnv, err := utils.EnsureRequiredEnvVars("pull")
	if err != nil {
		return err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	localDBName := localEnv["DB_DATABASE"]
	if localDBName == "" {
		localDBName = remoteDBName
	}
	localDBUser := cfg.DBUser
	localDBPassword := cfg.DbPassword

	fmt.Println("Importing database dump locally...")

	if err := runMySQLCommand(localDBUser, localDBPassword, "", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", localDBName)); err != nil {
		return fmt.Errorf("error creating local database: %v", err)
	}

	if err := runMySQLCommand(localDBUser, localDBPassword, localDBName, fmt.Sprintf("source %s", localPath)); err != nil {
		return fmt.Errorf("error importing database dump locally: %v", err)
	}

	if err := utils.RunRemoteCommand(env["PULL_SSH_USER"], env["PULL_HOST"], fmt.Sprintf("rm %s", dumpFile)); err != nil {
		return fmt.Errorf("error deleting remote database dump: %v", err)
	}

	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("error deleting local database dump: %v", err)
	}

	fmt.Println("Database successfully pulled!")

	if _, exists := localEnv["DB_DATABASE"]; !exists {
		fmt.Println("Adding database credentials to local .env file...")
		envManager := utils.NewEnvManager(".env")

		if err := envManager.SetVar("DB_DATABASE", localDBName); err != nil {
			return fmt.Errorf("error adding DB_DATABASE to .env: %v", err)
		}
		if err := envManager.SetVar("DB_USERNAME", localDBUser); err != nil {
			return fmt.Errorf("error adding DB_USERNAME to .env: %v", err)
		}
		if err := envManager.SetVar("DB_PASSWORD", localDBPassword); err != nil {
			return fmt.Errorf("error adding DB_PASSWORD to .env: %v", err)
		}

		fmt.Println("Database credentials added to local .env file.")
	} else {
		fmt.Println("Database credentials already exist in local .env file.")
	}
	return nil
}

func runMySQLCommand(user, password, dbName, query string) error {
	args := []string{"-u", user}
	if password != "" {
		args = append(args, "-p"+password)
	}
	if dbName != "" {
		args = append(args, dbName)
	}
	args = append(args, "-e", query)
	return utils.RunCommand("mysql", args...)
}

func getRemoteEnvValue(env map[string]string, remoteEnvPath, key, context string) (string, error) {
	var sshUserKey, sshHostKey string

	if context == "push" {
		sshUserKey = "PUSH_SSH_USER"
		sshHostKey = "PUSH_HOST"
	} else if context == "pull" {
		sshUserKey = "PULL_SSH_USER"
		sshHostKey = "PULL_HOST"
	} else {
		return "", fmt.Errorf("invalid context: %s", context)
	}

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", env[sshUserKey], env[sshHostKey]),
		fmt.Sprintf("grep %s %s | cut -d '=' -f 2", key, remoteEnvPath))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error fetching remote env value for %s: %v", key, err)
	}
	output = []byte(strings.Trim(strings.TrimSpace(string(output)), `"'`))
	return strings.TrimSpace(string(output)), nil
}
