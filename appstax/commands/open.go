package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"github.com/skratchdot/open-golang/open"
)

func DoOpen(c *cli.Context) {
	useOptions(c)
	dest := c.Args().First()
	log.Debugf("Open destination: %s", dest)
	url := ""
	message := ""

	switch dest {

	default: fallthrough
	case "deployed":
		app, err := account.GetCurrentApp()
		if err != nil {
			message = "Sorry, could not find a deployed app to open."
		}
		url = "http://" + app.HostingSubdomain + ".appstax.io"
		break

	case "admin":
		url = "http://appstax.com/admin/#/dashboard"
		break

	case "local":
		url = "http://localhost:9000/"
		break
	}

	if url != "" {
		term.Printf("Opening %s in your browser.\n", url)
		err := open.Start(url)
		if err != nil {
			message = "Ooops! Something went wrong."
		}
	}

	if message != "" {
		term.Section()
		term.Println(message)
	}
}
