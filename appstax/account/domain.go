package account

import (
	"sort"
	"strings"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Apps      []App  `json:"apps"`
}

type App struct {
	AppID            string       `json:"appID,omitempty"`
	AppKey           string       `json:"appKey,omitempty"`
	AppName          string       `json:"appName"`
	AppDescription   string       `json:"appDescription"`
	AccountID        string       `json:"accountId"`
	PaymentPlan      string       `json:"paymentPlan"`
	HostingSubdomain string       `json:"hostingSubdomain,omitempty"`
	Collections      []Collection `json:"collections,omitempty"`
}

type Collection struct {
	CollectionID   string                 `json:"collectionId,omitempty"`
	AppID          string                 `json:"appId"`
	AccountID      string                 `json:"accountId"`
	CollectionName string                 `json:"collectionName"`
	Schema         map[string]interface{} `json:"schema"`
}

func (app App) CollectionNames() []string {
	names := make([]string, 0)
	for _, collection := range app.Collections {
		names = append(names, collection.CollectionName)
	}
	return names
}

func (coll Collection) SortedColumnNames() []string {
	names := make([]string, 0)
	for k, _ := range coll.Schema["properties"].(map[string]interface{}) {
		names = append(names, k)
	}
	sorted     := make([]string, 0)
	sys := make([]string, 0)
	dev := make([]string, 0)
	for _, name := range names {
		if strings.HasPrefix(name, "sys") {
			sys = append(sys, name)
		} else {
			dev = append(dev, name)
		}
	}
	sort.Strings(sys)
	sort.Strings(dev)
	sorted = append(sorted, dev...)
	sorted = append(sorted, sys...)
	return sorted
}