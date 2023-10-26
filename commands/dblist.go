package commands

import (
	"database/sql"
	"fmt"
	"log"
	"mo/config" // if config is in a separate package

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

func ListDatabases(c *cli.Context) error {
	cfg, err := config.LoadConfig()

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/", cfg.DBUser, cfg.DbHost, cfg.DbPort)

	// fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		log.Fatalf("Error querying for databases: %v", err)
		return err
	}
	defer rows.Close()

	var databaseName string
	fmt.Println("----------------------------")
	fmt.Println("|      Database List       |")
	fmt.Println("----------------------------")
	for rows.Next() {
		err := rows.Scan(&databaseName)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
			return nil
		}
		fmt.Printf("| %-24s |\n", databaseName)
	}
	fmt.Println("----------------------------")
	return nil
}
