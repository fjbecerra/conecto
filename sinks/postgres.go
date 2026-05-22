package sinks

import (
	"conecto/core/codecs"
	"conecto/sinks/base/db"
	"fmt"
	"strings"
)

type Postgres struct{
	Codec codecs.Codec
	Upsert bool
}

func (a *Postgres) BuildQuery(schema db.Schema ,batchSize int) string {
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
	`,
		schema.Table,
		strings.Join(append(schema.Columns, metadataColumns...), ","),
		strings.Join(rows, ","),
	)
	if a.Upsert {
		query = query + fmt.Sprintf("ON CONFLICT (%s) DO NOTHING", strings.Join([]string{"__event_id"}, ""))
	}
	return query
}

func (a *Postgres) Decode(b []byte) (map[string]interface{}, error) {
	return a.Codec.Decode(b)
}