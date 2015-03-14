package session

import (
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
)

func WriteSessionID(sessionID string) {
	path := dir() + "/session"
	err := ioutil.WriteFile(path, []byte(sessionID), 0644)
	fail.Handle(err)
	log.Debugf("Wrote session id to %s", path)
}

func ReadSessionID() string {
	sessionID, _ := ioutil.ReadFile(dir() + "/session")
	if len(sessionID) == 0 {
		log.Debugf("Session id is empty")
	} else {
		log.Debugf("Read session id starting with %s", string(sessionID)[0:5]+"...")
	}
	return string(sessionID)
}

func WriteUserID(userID string) {
	path := dir() + "/user"
	err := ioutil.WriteFile(path, []byte(userID), 0644)
	fail.Handle(err)
	log.Debugf("Wrote user id to %s", path)
}

func ReadUserID() string {
	user, err := ioutil.ReadFile(dir() + "/user")
	fail.Handle(err)
	return string(user)
}

func WriteAccountID(accountID string) {
	path := dir() + "/account"
	err := ioutil.WriteFile(path, []byte(accountID), 0644)
	fail.Handle(err)
	log.Debugf("Wrote account id to %s", path)
}

func ReadAccountID() string {
	acc, err := ioutil.ReadFile(dir() + "/account")
	fail.Handle(err)
	return string(acc)
}

func Delete() {
	os.RemoveAll(dir())
}

func dir() string {
	home, err := homedir.Dir()
	fail.Handle(err)
	dir := home + "/.appstax/session/"
	err = os.MkdirAll(dir, 0700)
	fail.Handle(err)
	return dir
}
