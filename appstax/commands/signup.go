package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/term"
	"fmt"
	"github.com/codegangsta/cli"
)

func DoSignup(c *cli.Context) {
	useOptions(c)
	term.Println("By signing up you agree to our terms of service:")
	term.Println("https://appstax.com/admin/#/tos")
	term.Section()
	for {
		firstName := term.GetString("First name")
		lastName := term.GetString("Last name")
		email := term.GetString("Email")
		password := term.GetPassword("Password")
		sessionID, userID, accountID, err := account.Signup(firstName, lastName, email, password)
		if err != nil {
			term.Section()
			term.Println(err.Error())
		} else {
			writeSession(sessionID, userID, accountID)
			term.Section()
			term.Println(fmt.Sprintf("Account created for '%s'", email))
			return
		}
	}
}
