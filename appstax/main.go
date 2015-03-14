package main

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/apiclient"
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/hosting"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/session"
	"appstax-cli/appstax/term"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func main() {
	setupSignals()
	commands := setupCommands()
	log.Infof("Command to execute: %s", strings.Join(os.Args, " "))
	term.Section()
	commands.Run(os.Args)
	term.Section()
}

func setupCommands() *cli.App {
	app := cli.NewApp()
	app.Name = "appstax"
	app.Usage = "command line interface for appstax.com"
	app.Version = "0.9.0"

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
			Name:   "login",
			Usage:  "Login",
			Action: DoLogin,
			Flags:  app.Flags,
		},
		{
			Name:   "init",
			Usage:  "Initialize current directory as an appstax app",
			Action: DoInit,
			Flags:  app.Flags,
		},
		{
			Name:   "info",
			Usage:  "Info about app configured in current directory",
			Action: DoInfo,
			Flags:  app.Flags,
		},
		{
			Name:   "deploy",
			Usage:  "Deploy local files to <yourapp>.appstax.io",
			Action: DoDeploy,
			Flags:  app.Flags,
		},
		{
			Name:   "open",
			Usage:  "Open your browser to the specified destination",
			Action: DoOpen,
			Flags:  app.Flags,
		},
		{
			Name:   "logout",
			Usage:  "Logout of current app session",
			Action: DoLogout,
			Flags:  app.Flags,
		},
		{
			Name:   "serve",
			Usage:  "Run development http server",
			Action: DoServe,
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

func DoInit(c *cli.Context) {
	useOptions(c)
	loginIfNeeded()
	app := selectApp()
	pub := selectPublicDir()
	writeConfig(app, pub)
	selectSubdomain(app.AppID)
	createPublicDir()
	term.Section()
	term.Println("Now run 'appstax deploy' when you are ready to upload your public files.")
}

func DoInfo(c *cli.Context) {
	useOptions(c)
	loginIfNeeded()
	if !config.Exists() {
		term.Println("No app configured in current directory. (Missing appstax.conf)")
	} else {
		app, err := account.GetCurrentApp()
		fail.Handle(err)
		term.Println("App name:    " + app.AppName)
		term.Println("Description: " + app.AppDescription)
		term.Println("App key:     " + app.AppKey)
		term.Println("Hosting:     " + account.FormatHostingUrl(app))
		term.Section()
	}
	user, err := account.GetUser()
	fail.Handle(err)
	term.Printf("Logged in as %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
}

func DoLogin(c *cli.Context) {
	useOptions(c)
	login()
}

func DoDeploy(c *cli.Context) {
	useOptions(c)
	if !config.Exists() {
		term.Println("Can't find appstax.conf. Run 'appstax init' to initialize before deploying.")
		return
	}
	loginIfNeeded()
	dir := config.Read().PublicDir
	term.Println("Packaging files for upload...")
	archive, bytes := hosting.PrepareArchive(dir)
	term.Printf("Uploading %.2f MB...\n", float64(bytes)/(1024.0*1024.0))
	progress := term.ShowProgressBar(bytes)
	err := hosting.UploadStatic(archive, progress)
	if err != nil {
		term.Section()
		term.Println(err.Error())
	} else {
		progress.Finish()
		term.Section()
		term.Println("Deploy completed!")
	}
}

func DoLogout(c *cli.Context) {
	session.Delete()
	term.Println("Logged out.")
}

func DoOpen(c *cli.Context) {
	useOptions(c)
	dest := c.Args().First()
	log.Debugf("Open destination: %s", dest)
	url := ""
	message := ""

	switch dest {

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

func DoServe(c *cli.Context) {
	useOptions(c)
	publicDir := config.Read().PublicDir
	term.Println("Serving your public directory at http://localhost:9000/")
	term.Println("Press Ctrl-C to stop.")
	term.Section()

	http.Handle("/", http.FileServer(http.Dir(publicDir)))
	http.ListenAndServe(":9000", nil)
}

func useOptions(c *cli.Context) {
	log.SetStdoutEnabled(c.GlobalBool("verbose") || c.Bool("verbose"))
	fail.SetPanicMode(c.GlobalBool("verbose") || c.Bool("verbose"))
	apiclient.SetBaseUrl(c.GlobalString("baseurl") + c.String("baseurl"))
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

func selectApp() account.App {
	apps, _ := account.GetUserApps()
	selected := 0
	if len(apps) > 1 {
		term.Section()
		term.Println("Choose which app app to configure:")
		for i, app := range apps {
			term.Printf("  %d) %s\n", i+1, app.AppName)
		}
		term.Section()
		selected = -1 + term.GetInt(fmt.Sprintf("Please select app (1-%d)", len(apps)))
	}
	return apps[selected]
}

func selectPublicDir() string {
	dir := term.GetString("Select deployable directory [default: ./public]")
	if dir == "" {
		dir = "./public"
	}
	return dir
}

func selectSubdomain(appID string) {
	app, _ := account.GetAppByID(appID)
	log.Debugf("Subdomain app: %v", app)
	for {
		app.HostingSubdomain = term.GetString("Choose a *.appstax.io subdomain")
		err1 := account.SaveApp(app)
		if err1 != nil {
			term.Println(err1.Error())
		}
		err2 := account.AddCorsOrigin(appID, fmt.Sprintf("http://%s.appstax.io", app.HostingSubdomain))
		if err2 != nil {
			term.Println(err2.Error())
		}
		if err1 == nil && err2 == nil {
			term.Printf("Successfully configured %s.appstax.io\n", app.HostingSubdomain)
			return
		}
	}
}

func createPublicDir() {
	dir := config.Read().PublicDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		fail.Handle(err)
		ioutil.WriteFile(dir+"/index.html", []byte("<!DOCTYPE html>\n<h1>Your hosting is up and running</h1>\n<p>Now go make something amazing!</p>"), 0644)
		log.Debugf("Created public directory '%s'", dir)
	} else {
		log.Debugf("Not creating public directory. '%s' already exists.", dir)
	}
}

func writeConfig(app account.App, publicDir string) {
	config.Write(map[string]string{"AppKey": app.AppKey, "PublicDir": publicDir})
}

func writeSession(sessionID string, userID string, accountID string) {
	session.WriteSessionID(sessionID)
	session.WriteUserID(userID)
	session.WriteAccountID(accountID)
}
