package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"fmt"
	"strconv"
)

func DoFind(c *cli.Context) {
	useOptions(c)

	args := c.Args()
	if len(args) == 0 {
		term.Println("Too few arguments")
		return
	}

	collection, err := account.GetCollectionByName(args[0])
	if err != nil {
		term.Println(err.Error())
		return
	}

	filter := ""
	if len(args) > 1 {
		filter = args[1]
	}

	objects, err := account.GetObjects(collection.CollectionName, filter)
	if err != nil {
		term.Println(err.Error())
		return
	}

	columns := collection.SortedColumnNames()
	rows := objectsAsStringTable(columns, objects)
	term.PrintTable(columns, rows)
}

func objectsAsStringTable(columns []string, objects []map[string]interface{}) ([][]string) {
	rows  := make([][]string, 0)

	for _, object := range objects {
		row := make([]string, 0)
		for _, col := range columns {
			row = append(row, propertyAsString(object[col]))
		}
		rows = append(rows, row)
	}
	return rows
}

func propertyAsString(property interface{}) string {
	if stringValue, ok := property.(string); ok {
		return stringValue
	}
	if floatValue, ok := property.(float64); ok {
		return strconv.FormatFloat(floatValue, 'f', -1, 64)
	}
	if mapValue, ok := property.(map[string]interface{}); ok {
		dataType, ok := mapValue["sysDatatype"].(string)
		if ok && dataType == "file" {
			if filename, ok := mapValue["filename"].(string); ok {
				return filename
			}
		} else if ok && dataType == "relation" {
			objects, _    := mapValue["sysObjects"].([]interface{})
			collection, _ := mapValue["sysCollection"].(string)
			if collection == "" {
				collection = "objects"
			}
			return fmt.Sprintf("(%d %s)", len(objects), collection)
		}
	}
	return ""
}
