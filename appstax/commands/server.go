package commands

import (
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/hosting"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
)

func DoServer(c *cli.Context) {
	useOptions(c)
	if !config.Exists() {
		term.Println("Can't find appstax.conf. Run 'appstax init' to initialize before deploying.")
		return
	}
	loginIfNeeded()
	
	args := c.Args()
	if len(args) == 0 {
		term.Println("Too few arguments. Usage: appstax server create|delete")
		return
	}

	operation := args[0]
	switch operation {
	case "create":
		err := hosting.CreateServer()
		if err == nil {
			term.Println("Server created successfully!")
		} else {
			term.Println("Error creating server:")
			term.Println(err.Error())
		}
	case "delete":
		err := hosting.DeleteServer()
		if err == nil {
			term.Println("Server deleted!")
		} else {
			term.Println("Error deleting server:")
			term.Println(err.Error())
		}
	default:
		term.Printf("Unknown server operation '%s'\n", operation)
	}
}
