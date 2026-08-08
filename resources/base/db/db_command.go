package db

import "conecto/core/commands"


type SQLCommand struct {
    Query  string
    Values []interface{}
}

func (c *SQLCommand) Kind() string {
    return "sql"
}


func New(query string, values []interface{}) commands.Command {
	return &SQLCommand{
		Query: query,
		Values: values,
	}
}