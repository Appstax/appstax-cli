package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/session"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"errors"
	"strings"
)

func DoCollection(c *cli.Context) {
	useOptions(c)

	args := c.Args()
	if len(args) == 0 {
		term.Println("Too few arguments")
		return
	}

	name := args[0]
	collection, err := account.GetCollectionByName(name)
	if err == nil {
		term.Printf("Collection %s exists:\n", collection.CollectionName)
		if len(args) > 1 {
			term.Println("(Column arguments ignored)")
		}
		term.Section()
		showCollectionInfo(collection)
		return
	}

	schema, err := schemaWithColumns(args[1:])
	if err != nil {
		term.Println(err.Error())
		return;
	}

	createCollection(name, schema)
}

func showCollectionInfo(collection account.Collection) {
	properties := collection.Schema["properties"].(map[string]interface{})
	for pname, property := range properties {
		if !strings.HasPrefix(pname, "sys") {
			ptype := typeForSchemaProperty(property.(map[string]interface{}))
			term.Printf("%s:%s\n", pname, ptype)
		}
	}
	for pname, property := range properties {
		if strings.HasPrefix(pname, "sys") {
			ptype := typeForSchemaProperty(property.(map[string]interface{}))
			term.Printf("%s:%s\n", pname, ptype)
		}
	}
}

func typeForSchemaProperty(property map[string]interface{}) string {
	ptype := property["type"].(string)
	if ptype == "object" {
		details := property["properties"].(map[string]interface{})
		sysDatatype := details["sysDatatype"].(map[string]interface{})["pattern"].(string)
		if sysDatatype != "" {
			ptype = sysDatatype
		}
	}
	return ptype
}

func createCollection(name string, schema map[string]interface{}) {
	app, _ := account.GetCurrentApp()
	collection := account.Collection{
		CollectionName: name,
		AppID: app.AppID,
		AccountID: session.ReadAccountID(),
		Schema: schema,
	}

	collection, err := account.SaveNewCollection(collection)
	if err != nil {
		term.Println(err.Error())
	} else {
		term.Printf("Collection %s created successfully!\n", collection.CollectionName)
		term.Section()
		showCollectionInfo(collection)	
	}
}

func schemaWithColumns(columns []string) (map[string]interface{}, error) {
	schema := make(map[string]interface{})
	props := make(map[string]interface{})

	for _, column := range columns {
		parts := strings.Split(column, ":")
		if len(parts) != 2 {
			return schema, errors.New("Unrecognized column argument '" + column + "'. Expected format: <name>:<type>")
		}
		pname := parts[0]
		ptype := parts[1]
		if !isValidType(ptype) {
			return schema, errors.New("Unrecognized column type '" + ptype + "' in argument " + column)
		}
		props[pname] = makeSchemaProperty(ptype)
	}
	schema["properties"] = props
	return schema, nil
}

func makeSchemaProperty(ptype string) map[string]interface{} {
	prop := map[string]interface{}{}
	if ptype == "file" {
		prop["type"] = "object"
		prop["properties"] = map[string]interface{}{
			"sysDatatype": map[string]string{"type":"string", "pattern":"file"},
			"filename": map[string]string{"type":"string"},
			"url": map[string]string{"type":"string"},
		}
	} else {
		prop["type"] = ptype
	}
	return prop
}

func isValidType(t string) bool {
	switch t {
	case "string", "number", "file":
		return true
	default:
		return false
	}
}

