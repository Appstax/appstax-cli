package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/session"
	"appstax-cli/appstax/template"
	"appstax-cli/appstax/term"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func DoInit(c *cli.Context) {
	useOptions(c)
	loginIfNeeded()

	app, err := selectApp()
	if err != nil {
		return
	}

	tpl := selectTemplate()

	publicDir := "./public"
	serverDir := "./server"
	writeConfig(app, publicDir, serverDir)

	if !strings.HasPrefix(tpl.Name, "ios/") {
		createPublicDir()
	}
	createServerDir()

	if tpl.Name != "none" {
		term.Section()
		term.Println("Setting up template... ")
		template.Install(tpl)
		term.Println("Done.")
	}

	if !strings.HasPrefix(tpl.Name, "ios/") {
		term.Section()
		selectSubdomain(app.AppID, false)
	}

	term.Section()
	term.Println("All done!")
	if !strings.HasPrefix(tpl.Name, "ios/") {
		term.Println("Now run 'appstax deploy' when you are ready to upload your files.")
	}
}

func selectApp() (account.App, error) {
	apps, _ := account.GetUserApps()
	selected := -1
	if len(apps) == 0 {
		term.Section()
		term.Println("You have not created any apps yet! Create one now:")
	} else {
		term.Section()
		term.Println("Choose which app to configure or create a new one:")
		for i, app := range apps {
			term.Printf("  %d) %s\n", i+1, app.AppName)
		}
		term.Printf("  %d) Create a new app\n", len(apps)+1)
		term.Section()
		for selected < 0 || selected > len(apps) {
			selected = -1 + term.GetInt(fmt.Sprintf("Please select (1-%d)", len(apps)+1))
		}
	}
	if selected >= 0 && selected < len(apps) {
		return apps[selected], nil
	} else {
		return createApp()
	}
}

func createApp() (account.App, error) {
	app := account.App{}
	term.Section()
	app.AppName = term.GetString("App name")
	app.AppDescription = term.GetString("Description")
	app.AccountID = session.ReadAccountID()
	app.PaymentPlan = "PROTOTYPE"

	app, err := account.SaveNewApp(app)
	if err != nil {
		term.Section()
		term.Println(err.Error())
	} else {
		term.Section()
		term.Printf("Successfully created app '%s'\n", app.AppName)
	}

	return app, err
}

func selectTemplate() template.Template {
	templates := template.All()

	term.Section()
	term.Println("Choose a template for you app:")
	for i, template := range templates {
		term.Printf("  %d) %s\n", i+1, template.Label)
	}

	term.Section()
	for {
		selected := -1 + term.GetInt(fmt.Sprintf("Please select (1-%d)", len(templates)))
		if selected >= 0 && selected < len(templates) {
			return templates[selected]
		}
	}
}

func selectPublicDir() string {
	dir := term.GetString("Select deployable directory [default: ./public]")
	if dir == "" {
		dir = "./public"
	}
	return dir
}

func selectSubdomainIfNeeded() {
	app, err := account.GetCurrentApp()
	fail.Handle(err)
	if app.HostingSubdomain == "" {
		selectSubdomain(app.AppID, true)
	}
}

func selectSubdomain(appID string, required bool) {
	app, _ := account.GetAppByID(appID)
	log.Debugf("Subdomain app: %v", app)
	for {
		if required {
			app.HostingSubdomain = term.GetString("Choose a *.appstax.io subdomain for hosting")
		} else {
			term.Println("Choose a *.appstax.io subdomain for hosting:")
			term.Println("(Leave this blank if you wish to decide this later.)")
			app.HostingSubdomain = term.GetString("Subdomain")
		}

		if app.HostingSubdomain != "" {
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
		} else if !required {
			return
		}
	}
}

func createPublicDir() {
	dir := config.Read().PublicDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		fail.Handle(err)
		log.Debugf("Created public directory '%s'", dir)
	} else {
		log.Debugf("Not creating public directory. '%s' already exists.", dir)
	}
}

func createServerDir() {
	dir := config.Read().ServerDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		fail.Handle(err)
		log.Debugf("Created server directory '%s'", dir)
	} else {
		log.Debugf("Not creating server directory. '%s' already exists.", dir)
	}
}

