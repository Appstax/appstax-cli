package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/session"
	"appstax-cli/appstax/term"
	"fmt"
	"github.com/codegangsta/cli"
)

func DoLogin(c *cli.Context) {
	useOptions(c)
	login()
}

func loginIfNeeded() {
	if session.ReadSessionID() == "" {
		term.Section()
		term.Println("Please log in:")
		login()
	}
}

func login() {
	for {
		email := term.GetString("Email")
		password := term.GetPassword("Password")
		sessionID, userID, accountID, err := account.Login(email, password)
		if err != nil {
			term.Section()
			term.Println(err.Error())
		} else {
			writeSession(sessionID, userID, accountID)
			term.Section()
			term.Println(fmt.Sprintf("Successfully logged in as '%s'", email))
			return
		}
	}
}
