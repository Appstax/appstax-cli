package commands

import (
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/hosting"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
)

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
