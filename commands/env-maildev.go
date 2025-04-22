package commands

import (
	"fmt"
	"log"
	"mo/utils"

	"github.com/urfave/cli/v2"
)

func EnvMailDev(cliContext *cli.Context) error {
	envManager := utils.NewEnvManager(".env")

	if err := envManager.SetVar("MAIL_MAILER", "smtp"); err != nil {
		log.Fatal(err)
	}
	if err := envManager.SetVar("MAIL_HOST", "127.0.0.1"); err != nil {
		log.Fatal(err)
	}
	if err := envManager.SetVar("MAIL_PORT", "2525"); err != nil {
		log.Fatal(err)
	}
	if err := envManager.SetVar("MAIL_USERNAME", ""); err != nil {
		log.Fatal(err)
	}
	if err := envManager.SetVar("MAIL_PASSWORD", ""); err != nil {
		log.Fatal(err)
	}
	if err := envManager.SetVar("MAIL_ENCRYPTION", ""); err != nil {
		log.Fatal(err)
	}
	if err := envManager.SetVar("MAIL_FROM_ADDRESS", ""); err != nil {
		log.Fatal(err)
	}

	fmt.Println(".env file updated with MailDev settings.")

	return nil
}
