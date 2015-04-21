package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/apiclient"
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/session"
	"github.com/codegangsta/cli"
)

func useOptions(c *cli.Context) {
	log.SetStdoutEnabled(c.GlobalBool("verbose") || c.Bool("verbose"))
	fail.SetPanicMode(c.GlobalBool("verbose") || c.Bool("verbose"))
	apiclient.SetBaseUrl(c.GlobalString("baseurl") + c.String("baseurl"))
}

func writeConfig(app account.App, publicDir string) {
	config.Write(map[string]string{"AppKey": app.AppKey, "PublicDir": publicDir})
}

func writeSession(sessionID string, userID string, accountID string) {
	session.WriteSessionID(sessionID)
	session.WriteUserID(userID)
	session.WriteAccountID(accountID)
}
