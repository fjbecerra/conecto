package db

import (
	"conecto/core"
	"conecto/core/commands"
	"context"
	"fmt"
)

type Metadata struct{
	Column string
	Unique bool
}

type Schema struct {
	Table   string
	Columns []string
	Metadata []Metadata
}

type Rdbs struct {
	Schema  Schema
	Adapter Adapter
}

func (rdbs *Rdbs) WriteBatch(context context.Context, ID string, batch [] core.Event) ([]commands.Command, error) {
	fmt.Println("SINK: writing batch size =", len(batch))

	rows := make([]map[string]interface{}, 0, len(batch))
	
	for _, ev := range batch {
		//payload
		rec, err := rdbs.Adapter.Decode(ev.Payload)
		if err != nil {
			return nil,err
		}
		//metadata
		for key, value := range ev.Meta{
			rec[key]=value
		}
		rec["__stream_id"]=ID
		rows = append(rows, rec)
	}


	return rdbs.insertBatch(rows)
}

func (rdbs *Rdbs) insertBatch(batch []map[string]interface{},)  ([]commands.Command, error){

	query := rdbs.Adapter.BuildQuery(
		rdbs.Schema,
		len(batch),
	)

	values := make([]interface{}, 0, len(batch)*len(rdbs.Schema.Columns)*len(rdbs.Schema.Metadata))
	
	for _, rec := range batch {
		for _, col := range rdbs.Schema.Columns {
			values = append(values, rec[col])
		}
		for _, col := range rdbs.Schema.Metadata {
			values = append(values, rec[col.Column])
		}
	}
	
	return []commands.Command{
        &SQLCommand{
            Query: query,
            Values: values,
        },
    }, nil
}

