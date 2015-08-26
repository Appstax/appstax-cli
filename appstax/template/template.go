package template

import (
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/download"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Template struct {
	Name            string
	Label           string
	Repository      string
	SourcePath      string
	DestinationPath string
	AppKeyInFile    string
}

func All() []Template {
	return []Template{
		Template{
			Name:            "js/basic",
			Label:           "JavaScript: Basic project",
			Repository:      "appstax-js",
			SourcePath:      "appstax-js/starterprojects/basic/",
			DestinationPath: "./public/",
			AppKeyInFile:    "./public/app.js"},
		Template{
			Name:            "js/angular",
			Label:           "JavaScript: Basic angular.js project",
			Repository:      "appstax-js",
			SourcePath:      "appstax-js/starterprojects/angular/",
			DestinationPath: "./",
			AppKeyInFile:    "./app/modules/app.js"},
		Template{
			Name:            "js/angular",
			Label:           "JavaScript: Full angular.js project",
			Repository:      "appstax-js",
			SourcePath:      "appstax-js/starterprojects/angular-full/",
			DestinationPath: "./",
			AppKeyInFile:    "./app/modules/app.js"},
		Template{
			Name:            "js/library",
			Label:           "JavaScript: Just download the appstax.js file",
			Repository:      "appstax-js",
			SourcePath:      "appstax-js/starterprojects/basic/appstax.js",
			DestinationPath: "./public/appstax.js",
			AppKeyInFile:    ""},
		Template{
			Name:            "ios/basic",
			Label:           "iOS: Basic project",
			Repository:      "appstax-ios",
			SourcePath:      "appstax-ios/StarterProjects/Basic/",
			DestinationPath: "./",
			AppKeyInFile:    "./StarterProject/AppDelegate.m"},
		Template{
			Name:            "none",
			Label:           "No template",
			Repository:      "",
			SourcePath:      "",
			DestinationPath: "",
			AppKeyInFile:    ""},
	}
}

func Install(template Template) {
	if template.Repository == "" {
		return
	}

	releasePath := download.DownloadLatestRelease(template.Repository)
	log.Debugf("Downloaded release path: %s", releasePath)

	sourcePath := filepath.Join(releasePath, template.SourcePath)
	copy(sourcePath, config.ResolvePath(template.DestinationPath))
	insertAppKey(template)
}

func insertAppKey(template Template) {
	if template.AppKeyInFile == "" {
		return
	}
	path := config.ResolvePath(template.AppKeyInFile)
	bytes, err := ioutil.ReadFile(path)
	fail.Handle(err)
	text := string(bytes)

	text = strings.Replace(text, "<<appstax-app-key>>", config.Read().AppKey, -1)
	err = ioutil.WriteFile(path, []byte(text), 0644)
	fail.Handle(err)
}

func copy(src string, dst string) {
	log.Debugf("Copy '%s' to '%s'", src, dst)
	filepath.Walk(src, func(srcPath string, fileInfo os.FileInfo, err error) error {
		fail.Handle(err)
		if !fileInfo.IsDir() {
			log.Debugf("Copying file from '%s'", srcPath)
			dstPath := filepath.Join(dst, srcPath[len(src):])
			log.Debugf("... to destination '%s'", dstPath)
			srcFile, err := os.Open(srcPath)
			fail.Handle(err)
			defer srcFile.Close()
			os.MkdirAll(filepath.Dir(dstPath), 0755)
			dstFile, err := os.Create(dstPath)
			fail.Handle(err)
			defer dstFile.Close()
			_, err = io.Copy(dstFile, srcFile)
			fail.Handle(err)
		}
		return nil
	})
}
