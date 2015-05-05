package main

import (
	"appstax-cli/appstax/commands"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/term"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/signal"
	"strings"
)

func main() {
	setupSignals()
	cliApp := setupCli()
	log.Infof("Command to execute: %s", strings.Join(os.Args, " "))
	term.Section()
	cliApp.Run(os.Args)
	term.Section()
}

func setupCli() *cli.App {
	app := cli.NewApp()
	app.Name = "appstax"
	app.Usage = "command line interface for appstax.com"
	app.Version = "1.0.2"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose output",
		},
		cli.StringFlag{
			Name:  "baseurl",
			Usage: "set api base url",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "signup",
			Usage:  "Create new account",
			Action: commands.DoSignup,
			Flags:  app.Flags,
		},
		{
			Name:   "login",
			Usage:  "Login",
			Action: commands.DoLogin,
			Flags:  app.Flags,
		},
		{
			Name:   "init",
			Usage:  "Initialize current directory as an appstax app",
			Action: commands.DoInit,
			Flags:  app.Flags,
		},
		{
			Name:   "info",
			Usage:  "Info about app configured in current directory",
			Action: commands.DoInfo,
			Flags:  app.Flags,
		},
		{
			Name:   "deploy",
			Usage:  "Deploy local files to <yourapp>.appstax.io",
			Action: commands.DoDeploy,
			Flags:  app.Flags,
		},
		{
			Name:   "open",
			Usage:  "Open your browser to the specified destination",
			Action: commands.DoOpen,
			Flags:  app.Flags,
		},
		{
			Name:   "logout",
			Usage:  "Logout of current app session",
			Action: commands.DoLogout,
			Flags:  app.Flags,
		},
		{
			Name:   "serve",
			Usage:  "Run development http server",
			Action: commands.DoServe,
			Flags:  app.Flags,
		},
	}

	return app
}

func setupSignals() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Print("\n")
			os.Exit(-1)
		}
	}()
}
