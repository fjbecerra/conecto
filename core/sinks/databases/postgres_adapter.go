package databases

import (
	"conecto/core/sinks/codecs"
	"fmt"
	"strings"
)

type PostgresAdapter struct{
	Codec codecs.Codec
}

func (a *PostgresAdapter) BuildUpsertQuery(schema Schema, upsert Upsert ,batchSize int) string {

	var rows []string
	arg := 1
	columns:= append(schema.Columns, schema.Metadata...)
	for i := 0; i < batchSize; i++ {
		row := make([]string, len(columns))
		for j := range columns {
			row[j] = fmt.Sprintf("$%d", arg)
			arg++
		}
		rows = append(rows, "("+strings.Join(row, ",")+")")
	}

	var updates []string
	for _, c := range columns {
		updates = append(updates, fmt.Sprintf("%s=EXCLUDED.%s", c, c))
	}

	query:= fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES %s
		ON CONFLICT (%s) DO UPDATE SET %s
	`,
		schema.Table,
		strings.Join(columns, ","),
		strings.Join(rows, ","),
		strings.Join(upsert.ConflictColumns, ","),
		strings.Join(updates, ","),
	)
	return query
}

func (a *PostgresAdapter) Decode(b []byte) (map[string]interface{}, error) {
	return a.Codec.Decode(b)
}