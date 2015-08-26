package config

import (
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AppKey     string `json:"appKey"`
	PublicDir  string `json:"publicDir"`
	ServerDir  string `json:"serverDir"`
	ApiBaseUrl string `json:"apiBaseUrl,omitempty"`
}

const fileName = "appstax.conf"

func Exists() bool {
	return FilePath() != ""
}

func FilePath() (string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		return ""
	}
	tried := make([]string, 0)
	for {
		path := filepath.Join(dir, fileName)
		tried = append(tried, path)
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
		if len(dir) <= 3 {
			log.Warnf("Cound not find appstax.conf. Searched paths: %s", strings.Join(tried, ", "))
			return ""
		}
		dir = filepath.Dir(dir)
	}
}

func RootDir() string {
	return filepath.Dir(FilePath())
}

func ResolvePath(path string) string {
	return filepath.Join(RootDir(), path)
}

func Write(values map[string]string) {
	config := Read()
	config.AppKey = values["AppKey"]
	config.PublicDir = values["PublicDir"]
	config.ServerDir = values["ServerDir"]
	encoded, err := json.MarshalIndent(config, "", "    ")
	fail.Handle(err)
	ioutil.WriteFile(fileName, encoded, 0644)
	log.Debugf("Wrote config file: %s", encoded)
}

func Read() Config {
	var config Config
	dat, err := ioutil.ReadFile(FilePath())
	if err != nil {
		log.Debugf("Could not find appstax.conf")
	} else {
		err = json.Unmarshal(dat, &config)
		fail.Handle(err)
	}
	insertDefaults(&config)
	return config
}

func insertDefaults(config *Config) {
	if config.PublicDir == "" {
		config.PublicDir = "./public"
	}
	if config.ServerDir == "" {
		config.ServerDir = "./server"
	}
}
