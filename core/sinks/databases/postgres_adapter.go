package databases

import (
	"conecto/core/sinks/codecs"
	"fmt"
	"strings"
)

type PostgresAdapter struct{
	Codec codecs.Codec
}

func (a *PostgresAdapter) BuildUpsertQuery(schema Schema ,batchSize int) string {
	//schema
	var rows []string
	arg := 1
	for i := 0; i < batchSize; i++ {
		columns := make([]string, len(schema.Columns))
		for j := range schema.Columns {
			columns[j] = fmt.Sprintf("$%d", arg)
			arg++
		}
		metadata := make([]string, len(schema.Metadata))
		for j := range schema.Metadata {
			metadata[j] = fmt.Sprintf("$%d", arg)
			arg++
		}
		row:= append(columns, metadata...)
		rows = append(rows, "("+strings.Join(row, ",")+")")
	}

	metadataColumns := []string{}
	uniqueColumns := []string{}
	
	for _, c := range  schema.Metadata {
		metadataColumns = append(metadataColumns, c.Column)
		if c.Unique {
			uniqueColumns = append(uniqueColumns, c.Column)
		}
	}
	
	query:= fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES %s
		ON CONFLICT (%s) DO NOTHING
	`,
		schema.Table,
		strings.Join(append(schema.Columns, metadataColumns...), ","),
		strings.Join(rows, ","),
		strings.Join([]string{"__event_id","__pipeline_id"}, ","),
	)
	return query
}

func (a *PostgresAdapter) Decode(b []byte) (map[string]interface{}, error) {
	return a.Codec.Decode(b)
}