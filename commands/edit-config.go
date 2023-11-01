package commands

import (
	"mo/config"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func EditConfig(c *cli.Context) error {

	configPath, _ := config.ConfigPath()

	cfg, _ := config.LoadConfig()
	
	return exec.Command(cfg.Editor, configPath).Run()
}