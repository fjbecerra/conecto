package databases

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sinks/commands"
	"context"
	"fmt"
)


type Schema struct {
	Table   string
	Columns []string
	Metadata []string
}

type Upsert struct {
	ConflictColumns []string
}

type Rdbs struct {
	Schema  Schema
	Upsert  Upsert
	Adapter Adapter
}

func (rdbs *Rdbs) WriteBatch(ctx context.Context, batch [] core.Event) (sinks.Command, error) {
	fmt.Println("SINK: writing batch size =", len(batch))

	rows := make([]map[string]interface{}, 0, len(batch))

	for _, ev := range batch {

		rec, err := rdbs.Adapter.Decode(ev.Payload)
		if err != nil {
			return nil, err
		}

		rows = append(rows, rec)
	}

	return rdbs.insertBatch(ctx, rows),nil
}

func (rdbs *Rdbs) insertBatch(ctx context.Context, batch []map[string]interface{},) sinks.Command {

	query := rdbs.Adapter.BuildUpsertQuery(
		rdbs.Schema,
		rdbs.Upsert,
		len(batch),
	)

	columns:= append(rdbs.Schema.Columns, rdbs.Schema.Metadata...)
	values := make([]interface{}, 0, len(batch)*len(columns))
	
	

	for _, rec := range batch {
		for _, col := range columns {
			values = append(values, rec[col])
		}
	}
	
	return commands.New(query, values)
}

func (r *Rdbs) Open(ctx context.Context) error {
	fmt.Println("SINK: open")
	return nil
}

func (r *Rdbs) Close() error {
	defer r.Close()
	return nil
}


