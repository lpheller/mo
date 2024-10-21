package commands

import (
	"database/sql"
	"fmt"
	"log"
	"mo/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

func CreateDatabase(c *cli.Context) error {
	// Ensure database name is provided
	dbName := c.Args().First()
	if dbName == "" {
		return fmt.Errorf("missing database name")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf("%s@tcp(%s:%s)/", cfg.DBUser, cfg.DbHost, cfg.DbPort)

	// Open MySQL connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Use prepared statement for creating the database
	stmt := fmt.Sprintf("CREATE DATABASE `%s`", dbName)
	if _, err := db.Exec(stmt); err != nil {
		return fmt.Errorf("error creating database '%s': %w", dbName, err)
	}

	fmt.Printf("Database '%s' created successfully\n", dbName)
	return nil
}
