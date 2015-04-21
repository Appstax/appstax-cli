package download

import (
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"archive/zip"
	"errors"
	//"fmt"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadLatestRelease(repository string) string {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases("appstax", repository, nil)
	fail.Handle(err)

	if len(releases) == 0 {
		fail.Handle(errors.New("No releases found for '" + repository + "'"))
	}

	latest := releases[0]
	version := *latest.TagName
	log.Debugf("Latest release: %s", version)

	path := dir() + repository + "/" + version + "/"

	if _, err := os.Stat(path); err == nil {
		log.Debugf("Release already downloaded at %s", path)
		return path
	}

	asset, err := getAssetByName(latest, repository+".zip")
	fail.Handle(err)

	err = DownloadAndUnzip(*asset.BrowserDownloadURL, path)
	fail.Handle(err)

	return path
}

func DownloadAndUnzip(remoteUrl string, localPath string) error {
	log.Debugf("Downloading %s to %s", remoteUrl, localPath)

	tmp, err := download(remoteUrl)
	if err != nil {
		return err
	}

	log.Debugf("zip temporarily saved: %s", tmp)

	err = unzip(tmp, localPath)
	if err != nil {
		return err
	}

	return nil
}

func download(url string) (path string, err error) {
	tmp, err := ioutil.TempFile("", "")
	defer tmp.Close()
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		return "", err
	}

	return tmp.Name(), nil
}

func unzip(src string, dest string) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	os.MkdirAll(dest, 0755)
	for _, file := range zipReader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		path := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}
			file.Close()
		}
	}

	return nil
}

func getAssetByName(release github.RepositoryRelease, name string) (github.ReleaseAsset, error) {
	for _, asset := range release.Assets {
		log.Debugf("Available asset: '%s'", *asset.Name)
		if *asset.Name == name {
			return asset, nil
		}
	}
	return github.ReleaseAsset{}, errors.New("Asset '" + name + "' not found")
}

func dir() string {
	home, err := homedir.Dir()
	fail.Handle(err)
	dir := home + "/.appstax/download/"
	err = os.MkdirAll(dir, 0700)
	fail.Handle(err)
	return dir
}
