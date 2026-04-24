package rdbs

import (
	"fmt"
	"strings"
)

type PostgresAdapter struct{
}

func (a *PostgresAdapter) BuildUpsertQuery(schema Schema, upsert Upsert ,batchSize int) string {

	var rows []string
	arg := 1

	for i := 0; i < batchSize; i++ {
		row := make([]string, len(schema.Columns))
		for j := range schema.Columns {
			row[j] = fmt.Sprintf("$%d", arg)
			arg++
		}
		rows = append(rows, "("+strings.Join(row, ",")+")")
	}

	var updates []string
	for _, c := range schema.Columns {
		updates = append(updates, fmt.Sprintf("%s=EXCLUDED.%s", c, c))
	}

	return fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES %s
		ON CONFLICT (%s) DO UPDATE SET %s
	`,
		schema.Table,
		strings.Join(schema.Columns, ","),
		strings.Join(rows, ","),
		strings.Join(upsert.ConflictColumns, ","),
		strings.Join(updates, ","),
	)
}