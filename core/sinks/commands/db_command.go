package commands

import "conecto/core/sinks"

type SqlCommand struct{
	Query string
	Values []interface{}
}

func New(query string, values ...interface{}) sinks.Command {
	return SqlCommand{
		Query: query,
		Values: values,
	}
}

