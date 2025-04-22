package commands

import (
	"database/sql"
	"fmt"
	"log"
	"mo/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

func CreateDatabase(cliContext *cli.Context) error {
	dbName := cliContext.Args().First()
	if dbName == "" {
		return fmt.Errorf("missing database name")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/", cfg.DBUser, cfg.DbHost, cfg.DbPort)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	stmt := fmt.Sprintf("CREATE DATABASE `%s`", dbName)
	if _, err := db.Exec(stmt); err != nil {
		return fmt.Errorf("error creating database '%s': %w", dbName, err)
	}

	fmt.Printf("Database '%s' created successfully\n", dbName)
	return nil
}
