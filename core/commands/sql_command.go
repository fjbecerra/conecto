package commands


type SQLCommand struct {
    Query  string
    Values []any
}

func (c *SQLCommand) Kind() string {
    return "sql"
}
