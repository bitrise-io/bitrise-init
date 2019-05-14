package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/bitrise-io/bitrise-init/version"
	"github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

// Run ...
func Run() {
	// Parse cl
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Version)
	}

	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "Bitrise Init Tool"
	app.Version = version.VERSION
	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "ci",
			Usage:  "If true it indicates that we're used by another tool so don't require any user input!",
			EnvVar: "CI",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetEnableDebugLog(true)
		return nil
	}

	app.Commands = []cli.Command{
		versionCommand,
		configCommand,
		manualConfigCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%s", err)
		os.Exit(1)
	}
}
