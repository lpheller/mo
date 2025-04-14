package commands

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

	"mo/utils"
    "mo/config"

    "github.com/urfave/cli/v2"
)

func Pull(c *cli.Context) error {
    localEnv, err := loadLocalEnv()
    if (err != nil) {
        return err
    }

    // if no flags are set, print usage and return
    if !c.Bool("storage") && !c.Bool("database") {
        fmt.Println("No flags set. Use --storage or --database to pull data.")
        return nil
    }

    // Check if the --storage flag is set, and pull the storage folder if true.
    if c.Bool("storage") {
        if err := pullStorage(localEnv); err != nil {
            return err
        }
    }

    // Check if the --database flag is set, and pull the database if true.
    if c.Bool("database") {
        if err := pullDatabase(localEnv); err != nil {
            return err
        }
    }

    return nil
}

func loadLocalEnv() (map[string]string, error) {
    requiredKeys := []string{"REMOTE_SSH_USER", "REMOTE_IP", "REMOTE_PROJECT_DIR"}
    return utils.LoadEnv(".env", requiredKeys)
}

func pullStorage(env map[string]string) error {
    fmt.Println("Pulling storage folder...")

    dateOutput, err := exec.Command("date", "+%F").Output()
    if err != nil {
        return fmt.Errorf("error getting date: %v", err)
    }
    storageFile := fmt.Sprintf("storage_%s.tar.gz", strings.TrimSpace(string(dateOutput)))
    remoteCmd := fmt.Sprintf("cd %s && tar -cvzf /tmp/%s -C storage/app/public .", env["REMOTE_PROJECT_DIR"], storageFile)

    if err := runRemoteCommand(env, remoteCmd); err != nil {
        return fmt.Errorf("error compressing storage folder on remote: %v", err)
    }

    localPath := fmt.Sprintf("/tmp/%s", storageFile)
    remotePath := fmt.Sprintf("%s@%s:/tmp/%s", env["REMOTE_SSH_USER"], env["REMOTE_IP"], storageFile)

    if err := runCommand("scp", remotePath, localPath); err != nil {
        return fmt.Errorf("error downloading storage file: %v", err)
    }

    if err := runCommand("tar", "-xzf", localPath, "-C", "storage/app/public"); err != nil {
        return fmt.Errorf("error extracting storage file locally: %v", err)
    }

    if err := runRemoteCommand(env, fmt.Sprintf("rm /tmp/%s", storageFile)); err != nil {
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

    remoteEnvPath := fmt.Sprintf("%s/.env", env["REMOTE_PROJECT_DIR"])
    fmt.Println("Remote env path:", remoteEnvPath)
    remoteDBName := getRemoteEnvValue(env, remoteEnvPath, "DB_DATABASE")
    fmt.Println("Remote DB Name: %s\n", remoteDBName)
    
    remoteDBUser := getRemoteEnvValue(env, remoteEnvPath, "DB_USERNAME")
    fmt.Println("Remote DB User: %s\n", remoteDBUser)
    remoteDBPassword := getRemoteEnvValue(env, remoteEnvPath, "DB_PASSWORD")
    


    dumpFile := fmt.Sprintf("/tmp/%s-dump.sql", remoteDBName)
    remoteCmd := fmt.Sprintf("mysqldump -u %s -p%s %s > %s", remoteDBUser, remoteDBPassword, remoteDBName, dumpFile)

    fmt.Println("Remote command to create dump:", remoteCmd)

    if err := runRemoteCommand(env, remoteCmd); err != nil {
        return fmt.Errorf("error creating database dump on remote: %v", err)
    }

    localPath := fmt.Sprintf("%s-dump.sql", remoteDBName)
    remotePath := fmt.Sprintf("%s@%s:%s", env["REMOTE_SSH_USER"], env["REMOTE_IP"], dumpFile)

    if err := runCommand("scp", remotePath, localPath); err != nil {
        return fmt.Errorf("error downloading database dump: %v", err)
    }
    // fmt.Println("Local path for database dump:", localPathV)
    // fmt.Println("Remote path for database dump:", remotePath)

    localEnv, err := loadLocalEnv()
    if err != nil {
        return err
    }

    // Laden der Konfiguration
    cfg, err := config.LoadConfig()
    if err != nil {
        return fmt.Errorf("error loading config: %w", err)
    }

    localDBName := localEnv["DB_DATABASE"] // Beibehalten, da der DB-Name aus der .env-Datei kommt
    if localDBName == "" {
        localDBName = remoteDBName // Use the same name as remote if not set
    }
    localDBUser := cfg.DBUser             // Benutzername aus der Konfiguration
    localDBPassword := cfg.DbPassword     // Passwort aus der Konfiguration

    fmt.Println("Local DB Name:", localDBName)
    fmt.Println("Local DB User:", localDBUser)
    fmt.Println("Local DB Password:", localDBPassword)
    fmt.Println("Importing database dump locally...")
    
    // Check if the local database exists
    if localDBPassword != "" {
        if err := runCommand("mysql", "-u", localDBUser, "-p"+localDBPassword, "-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", localDBName)); err != nil {
            return fmt.Errorf("error creating local database: %v", err)
        }
    } else {
        if err := runCommand("mysql", "-u", localDBUser, "-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", localDBName)); err != nil {
            return fmt.Errorf("error creating local database: %v", err)
        }
    }

    if localDBPassword != "" {
        if err := utils.RunCommand("mysql", "-u", localDBUser, "-p"+localDBPassword, localDBName, "-e", fmt.Sprintf("source %s", localPath)); err != nil {
            return fmt.Errorf("error importing database dump locally: %v", err)
        }
    } else {
        if err := utils.RunCommand("mysql", "-u", localDBUser, localDBName, "-e", fmt.Sprintf("source %s", localPath)); err != nil {
            return fmt.Errorf("error importing database dump locally: %v", err)
        }
    }

    if err := runRemoteCommand(env, fmt.Sprintf("rm %s", dumpFile)); err != nil {
        return fmt.Errorf("error deleting remote database dump: %v", err)
    }

    if err := os.Remove(localPath); err != nil {
        return fmt.Errorf("error deleting local database dump: %v", err)
    }

    fmt.Println("Database successfully pulled!")

    // if the local env file does not contain the database credentials, add them
    if _, exists := localEnv["DB_DATABASE"]; !exists {
        fmt.Println("Adding database credentials to local .env file...")
        file, err := os.OpenFile(".env", os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            return fmt.Errorf("error opening .env file: %v", err)
        }
        defer file.Close()
        if _, err := file.WriteString(fmt.Sprintf("DB_DATABASE=%s\n", localDBName)); err != nil {   
            return fmt.Errorf("error writing to .env file: %v", err)
        }
        if _, err := file.WriteString(fmt.Sprintf("DB_USERNAME=%s\n", localDBUser)); err != nil {
            return fmt.Errorf("error writing to .env file: %v", err)
        }
        if _, err := file.WriteString(fmt.Sprintf("DB_PASSWORD=%s\n", localDBPassword)); err != nil {
            return fmt.Errorf("error writing to .env file: %v", err)
        }
        fmt.Println("Database credentials added to local .env file.")
    } else {
        fmt.Println("Database credentials already exist in local .env file.")
    }
    return nil
}

func runRemoteCommand(env map[string]string, cmd string) error {
    sshCmd := exec.Command("ssh", fmt.Sprintf("%s@%s", env["REMOTE_SSH_USER"], env["REMOTE_IP"]), cmd)
    sshCmd.Stdout = os.Stdout
    sshCmd.Stderr = os.Stderr
    return sshCmd.Run()
}

func runCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

func getRemoteEnvValue(env map[string]string, remoteEnvPath, key string) string {
    cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", env["REMOTE_SSH_USER"], env["REMOTE_IP"]),
        fmt.Sprintf("grep %s %s | cut -d '=' -f 2", key, remoteEnvPath))
    output, err := cmd.Output()
    if err != nil {
        log.Fatalf("error fetching remote env value for %s: %v", key, err)
    }
    // if the value is wrapped in quotes, remove them
    output = []byte(strings.Trim(strings.TrimSpace(string(output)), `'`))
    // if the value is wrapped in double quotes, remove them
    output = []byte(strings.Trim(strings.TrimSpace(string(output)), `"`))
    
    return strings.TrimSpace(string(output))
}