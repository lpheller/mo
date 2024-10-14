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
	// Check if the user passed a database name as the first argument
	if c.Args().Len() == 0 {
		log.Fatalf("You must provide a database name")
		return fmt.Errorf("missing database name")
	}

	dbName := c.Args().Get(0)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return err
	}

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/", cfg.DBUser, cfg.DbHost, cfg.DbPort)

	// Open MySQL connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return err
	}
	defer db.Close()

	// Create the new database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
		return err
	}

	fmt.Printf("Database '%s' created successfully\n", dbName)
	return nil
}
