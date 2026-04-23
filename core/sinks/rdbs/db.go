package rdbs

import (
	"conecto/core"
	"context"
	"database/sql"
)

type Field struct {
	Name string
	Default interface{}
}

type Schema struct {
	Fields []Field
} 

type Rdbs struct {
	DB      *sql.DB
	Table   string
	Schema  Schema
	Adapter Adapter
	BatchSize int
}


func (rdbs *Rdbs) Write(ctx context.Context, in <-chan core.Record) <- chan error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		batch := []map[string]interface{}{}
		batchSize := rdbs.BatchSize

		for rec := range in {
			batch = append(batch, rec)

			if len(batch) >= batchSize {
				if err := rdbs.insertBatch(ctx, batch); err != nil {
					errCh <- err
					return
				}
				batch = batch[:0]
			}
		}

		// flush
		if len(batch) > 0 {
			//todo retry
			if err := rdbs.insertBatch(ctx, batch); err != nil {
				errCh <- err
				return
			}
		}
	}()

	return errCh

}

func (rdbs *Rdbs) normalize(rec map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for _, f := range rdbs.Schema.Fields {
		val, ok := rec[f.Name]
		if !ok || val == nil {
			out[f.Name] = f.Default
		} else {
			out[f.Name] = val
		}
	}

	return out
}

func (rdbs *Rdbs) insertBatch(ctx context.Context, batch []map[string]interface{}) error {

	columns := make([]string, 0)
	for _, f := range rdbs.Schema.Fields {
		columns = append(columns, f.Name)
	}

	query := rdbs.Adapter.BuildInsert(rdbs.Table, columns, len(batch))

	args := []interface{}{}
	for _, rec := range batch {
		rec = rdbs.normalize(rec)
		for _, col := range columns {
			args = append(args, rec[col])
		}
	}

	_, err := rdbs.DB.ExecContext(ctx, query, args...)
	return err
}


