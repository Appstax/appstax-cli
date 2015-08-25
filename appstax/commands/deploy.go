package commands

import (
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/hosting"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"os"
)

func DoDeploy(c *cli.Context) {
	useOptions(c)
	if !config.Exists() {
		term.Println("Can't find appstax.conf. Run 'appstax init' to initialize before deploying.")
		return
	}
	loginIfNeeded()
	deployPublic()
	deployServer()
}

func deployPublic() {
	selectSubdomainIfNeeded()
	dir := config.Read().PublicDir
	term.Section()
	if !directoryExists(dir) {
		term.Println("Directory does not exist: "+dir)
		term.Println("Skipping public deploy")
		return
	}
	term.Println("Packaging public files for upload...")
	archive, bytes, _ := hosting.PrepareArchive(dir)
	term.Printf("Uploading %.2f MB...\n", float64(bytes)/(1024.0*1024.0))
	progress := term.ShowProgressBar(bytes)
	err := hosting.UploadStatic(archive, progress)
	if err != nil {
		progress.Finish()
		term.Section()
		term.Println("Error deploying public files: "+err.Error())
	} else {
		progress.Finish()
		term.Section()
		term.Println("Public deploy completed!")
	}
}

func deployServer() {
	dir := config.Read().ServerDir
	term.Section()
	if !directoryExists(dir) {
		term.Println("Directory does not exist: "+dir)
		term.Println("Skipping server deploy")
		return
	}
	term.Println("Packaging server files for upload...")
	archive, bytes, err := hosting.PrepareArchive(dir)
	if err != nil {
		term.Section()
		term.Println(err.Error())
		return
	}
	term.Printf("Uploading %.2f MB...\n", float64(bytes)/(1024.0*1024.0))
	progress := term.ShowProgressBar(bytes)

	err = hosting.UploadServer(archive, progress)
	if err != nil {
		progress.Finish()
		term.Section()
		term.Println("Error deploying server: "+err.Error())
	} else {
		progress.Finish()
		term.Section()
		term.Println("Server deploy completed!")
	}
}

func directoryExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil || !os.IsNotExist(err)
}