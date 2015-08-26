package commands

import (
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"net/http"
)

func DoServe(c *cli.Context) {
	useOptions(c)
	publicDir := config.ResolvePath(config.Read().PublicDir)
	term.Println("Serving your public directory at http://localhost:9000/")
	term.Println("Press Ctrl-C to stop.")
	term.Section()

	http.Handle("/", http.FileServer(http.Dir(publicDir)))
	http.ListenAndServe(":9000", nil)
}
