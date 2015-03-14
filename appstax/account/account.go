package account

import (
	"appstax-cli/appstax/apiclient"
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/session"
	"encoding/json"
	"errors"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Apps      []App  `json:"apps"`
}

type App struct {
	AppID            string `json:"appID"`
	AppKey           string `json:"appKey"`
	AppName          string `json:"appName"`
	AppDescription   string `json:"appDescription"`
	HostingSubdomain string `json:"hostingSubdomain"`
}

func Login(email string, password string) (sessionID string, userID string, accountID string, err error) {
	result, resp, err := apiclient.Post(apiclient.Url("/appstax/sessions"),
		map[string]interface{}{"email": email, "password": password})
	if err != nil {
		return "", "", "", err
	} else {
		resultMap := apiclient.ParseStringMap(result)
		log.Debugf("Login result: %v", resultMap)
		return resp.Header.Get("x-appstax-sessionid"), resultMap["userId"], resultMap["accountId"], nil
	}
}

func GetUser() (User, error) {
	userID := session.ReadUserID()
	result, _, _ := apiclient.Get(apiclient.Url("/appstax/users/" + userID))
	var user User
	err := json.Unmarshal(result, &user)
	return user, err
}

func GetUserApps() ([]App, error) {
	userID := session.ReadUserID()
	result, _, err := apiclient.Get(apiclient.Url("/appstax/users/" + userID))
	if err != nil {
		return nil, err
	}
	resultMap := apiclient.ParseMap(result)
	resultApps := resultMap["apps"].([]interface{})
	apps := make([]App, 0)
	for _, resultApp := range resultApps {
		app := resultApp.(map[string]interface{})
		apps = append(apps, App{
			AppID:          app["appId"].(string),
			AppKey:         app["appKey"].(string),
			AppName:        app["appName"].(string),
			AppDescription: app["appDescription"].(string),
		})
	}
	return apps, nil
}

func GetAppByID(appID string) (App, error) {
	result, _, _ := apiclient.Get(apiclient.Url("/appstax/apps/" + appID))
	var app App
	err := json.Unmarshal(result, &app)
	return app, err
}

func GetAppByKey(appKey string) (App, error) {
	apps, err := GetUserApps()
	fail.Handle(err)
	for _, app := range apps {
		if app.AppKey == appKey {
			return app, nil
		}
	}
	return App{}, errors.New("App not found")
}

func GetCurrentApp() (App, error) {
	appKey := config.Read().AppKey
	app, err := GetAppByKey(appKey)
	fail.Handle(err)
	app, err = GetAppByID(app.AppID)
	return app, err
}

func SaveApp(app App) error {
	_, _, err := apiclient.Put(apiclient.Url("/appstax/apps/"+app.AppID), app)
	return err
}

func AddCorsOrigin(appID string, origin string) error {
	origins := GetCorsOrigins(appID)
	log.Debugf("Existing CORS origins for app %s: %v", appID, origins)
	if -1 == indexOfString(origin, origins) {
		origins = append(origins, origin)
		_, _, err := apiclient.Put(apiclient.Url("/appstax/apps/"+appID+"/origins"), origins)
		fail.Handle(err)
		log.Debugf("Added new CORS origin: %v", origins)
	}
	return nil
}

func GetCorsOrigins(appID string) []string {
	result, _, err := apiclient.Get(apiclient.Url("/appstax/apps/" + appID + "/origins"))
	fail.Handle(err)
	var origins []string
	err = json.Unmarshal(result, &origins)
	fail.Handle(err)
	return origins
}

func FormatHostingUrl(app App) string {
	if app.HostingSubdomain == "" {
		return ""
	} else {
		return "http://" + app.HostingSubdomain + ".appstax.io"
	}
}

func indexOfString(needle string, haystack []string) int {
	for i, v := range haystack {
		if needle == v {
			return i
		}
	}
	return -1
}
