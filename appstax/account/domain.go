package account

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
