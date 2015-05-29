package commands

import (
	"appstax-cli/appstax/account"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"errors"
	"strings"
)

func DoRelation(c *cli.Context) {
	useOptions(c)

	args := c.Args()
	if len(args) < 2 {
		term.Println("Too few arguments")
		return
	}

	from := args[0]
	to   := args[1]
	relation, err := makeRelationFromArgs(from, to)
	if err != nil {
		term.Println(err.Error())
		return
	}

	_, err = account.SaveNewRelation(relation)
	if err != nil {
		term.Println(err.Error())
		return
	}

	term.Println("Succesfully added " + relation.Description())
}

func makeRelationFromArgs(from, to string) (account.Relation, error) {
	relation := account.Relation{}

	fromParts := strings.Split(from, ".")
	if len(fromParts) != 2 {
		return relation, errors.New("Unrecognized argument '" + from + "'. Expected format: <collection>:<property>")
	}

	fromCollectionName := fromParts[0]
	fromPropertyName   := fromParts[1]
	relationType := "single"
	if strings.HasSuffix(fromPropertyName, "[]") {
		relationType = "array"
		fromPropertyName = strings.Replace(fromPropertyName, "[]", "", 1)
	}

	relation.ToCollectionName = to
	relation.FromCollectionName = fromCollectionName
	relation.FromProperty = fromPropertyName
	relation.Type = relationType
	return relation, nil
}
