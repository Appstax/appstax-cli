package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"strings"
)

func DoInfo(c *cli.Context) {
	useOptions(c)
	loginIfNeeded()
	if !config.Exists() {
		term.Println("No app configured in current directory. (Missing appstax.conf)")
	} else {
		app, err := account.GetCurrentApp()
		if err != nil {
			term.Println("You don't have access to the currently selected app")
		} else {
			term.Println("App name:    " + app.AppName)
			term.Println("Description: " + app.AppDescription)
			term.Println("App key:     " + app.AppKey)
			term.Println("Collections: " + strings.Join(app.CollectionNames(), ", "))
			term.Println("Hosting:     " + account.FormatHostingUrl(app))
			term.Section()
		}
	}
	user, err := account.GetUser()
	fail.Handle(err)
	term.Printf("Logged in as %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
}
