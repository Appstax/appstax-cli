package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/template"
	"appstax-cli/appstax/term"
	"errors"
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

	pub := "./public"
	if tpl.Name == "none" {
		term.Section()
		pub = selectPublicDir()
	}
	writeConfig(app, pub)

	if !strings.HasPrefix(tpl.Name, "ios/") {
		createPublicDir()
	}

	if tpl.Name != "none" {
		term.Section()
		term.Println("Setting up template... ")
		template.Install(tpl)
		term.Println("Done.")
	}

	term.Section()
	selectSubdomain(app.AppID)

	term.Section()
	term.Println("All done!")
	term.Println("Now run 'appstax deploy' when you are ready to upload your public files.")
}

func selectApp() (account.App, error) {
	apps, _ := account.GetUserApps()
	selected := 0
	if len(apps) == 0 {
		term.Section()
		term.Println("You don't have access to any apps!")
		term.Println("Create one on appstax.com and come back here.")
		return account.App{}, errors.New("Account has no apps")
	} else if len(apps) > 1 {
		term.Section()
		term.Println("Choose which app to configure:")
		for i, app := range apps {
			term.Printf("  %d) %s\n", i+1, app.AppName)
		}
		term.Section()
		selected = -1 + term.GetInt(fmt.Sprintf("Please select (1-%d)", len(apps)))
	}
	return apps[selected], nil
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

func selectSubdomain(appID string) {
	app, _ := account.GetAppByID(appID)
	log.Debugf("Subdomain app: %v", app)
	for {
		app.HostingSubdomain = term.GetString("Choose a *.appstax.io subdomain for web hosting")
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
		log.Debugf("Created public directory '%s'", dir)
	} else {
		log.Debugf("Not creating public directory. '%s' already exists.", dir)
	}
}
