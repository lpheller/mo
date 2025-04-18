package commands

import (
	"fmt"
	"log"

	"mo/utils"

	"github.com/urfave/cli/v2"
)

func EnvSqlite(c *cli.Context) error {
	env := utils.NewEnvManager(".env")

	val, found, err := env.GetVar("DB_CONNECTION")
	if err != nil {
		log.Fatal(err)
	}
	if found {
		if val == "sqlite" {
			fmt.Println("DB_CONNECTION is already set to sqlite.")
			return nil
		}
	}

	err = env.SetVar("DB_CONNECTION", "sqlite")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB_CONNECTION set to sqlite.")

	return nil

}
