package commands

import (
	"appstax-cli/appstax/session"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
)

func DoLogout(c *cli.Context) {
	session.Delete()
	term.Println("Logged out.")
}
