package main

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "Blueprint"
	app.Copyright = "(c) 1996 Wercker Holding BV"
	app.Usage = "TiVo for VRML"

	//app.Version = version.Version
	//app.Compiled = version.CompiledAt

	app.Flags = []cli.Flag{
		debugFlag,
	}
	app.Commands = []cli.Command{
		//clientCommand,
		gatewayCommand,
		serverCommand,
	}

	app.Run(os.Args)
}

var debugFlag cli.BoolFlag = cli.BoolFlag{
	Name: "Debug",
}

// GlobalOptions are global
type GlobalOptions struct {
	Debug bool
}

func ParseGlobalOptions(c *cli.Context) (*GlobalOptions, error) {
	return &GlobalOptions{
		Debug: c.GlobalBool("debug"),
	}, nil
}

var ErrorExitCode = cli.NewExitError("", 1)
