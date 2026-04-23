package rdbs

import (
	"fmt"
	"strings"
)

type PostgresAdapter struct{}

func (a *PostgresAdapter) DriverName() string {
	return "postgres"
}

func (a *PostgresAdapter) Placeholder(n int) string {
	return fmt.Sprintf("$%d", n)
}

func (a *PostgresAdapter) BuildInsert(table string, columns []string, rows int) string {
	var parts []string
	arg := 1

	for i := 0; i < rows; i++ {
		row := []string{}
		for range columns {
			row = append(row, a.Placeholder(arg))
			arg++
		}
		parts = append(parts, "("+strings.Join(row, ",")+")")
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		table,
		strings.Join(columns, ","),
		strings.Join(parts, ","),
	)
}