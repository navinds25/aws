package app

import "github.com/urfave/cli"

type MainFlag struct {
	Help    bool
	Version bool
}

var MainFlagVal MainFlag

// Cli for all commandline arguments.
func Cli() *cli.App {
	app := cli.NewApp()
	cli.HelpFlag = cli.BoolFlag{
		Name:        "h",
		Destination: &MainFlagVal.Help,
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:        "v",
		Destination: &MainFlagVal.Version,
	}
	app.Name = "kms-pass"
	app.Usage = "The cloud password manager"
	app.Commands = []cli.Command{
		sftpConfCli,
	}
	app.Flags = mainCliFlags
	return app
}
