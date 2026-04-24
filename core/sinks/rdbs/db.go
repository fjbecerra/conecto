package rdbs

import (
	"conecto/core"
	"context"
	"database/sql"
	"fmt"
)


type Schema struct {
	Table   string
	Columns []string
}

type Upsert struct {
	ConflictColumns []string
}

type Rdbs struct {
	DB      *sql.DB	
	Schema  Schema
	Upsert  Upsert
	Adapter Adapter
	BatchSize int
}


func (rdbs *Rdbs) Write(ctx context.Context, in <-chan core.Record) <- chan error {
	fmt.Println("SINK: started")
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		tx, err := rdbs.DB.BeginTx(ctx, nil)
		if err != nil {
			errCh <- err
			return
		}

		batch := make([]map[string]interface{}, 0, rdbs.BatchSize)

		for rec := range in {
			batch = append(batch, rec)

			if len(batch) >= rdbs.BatchSize {
				if err := rdbs.insertBatchTx(ctx, tx, batch); err != nil {
					tx.Rollback()
					errCh <- err
					return
				}
				batch = batch[:0]
			}
		}		
		fmt.Println("SINK: finished reading input channel")

		// flush remaining
		if len(batch) > 0 {
			if err := rdbs.insertBatchTx(ctx, tx, batch); err != nil {
				tx.Rollback()
				errCh <- err
				return
			}
		}
		
		fmt.Println("SINK: committing transaction")
		if err := tx.Commit(); err != nil {
			errCh <- err
			return
		}
		fmt.Println("SINK: commit done")
		
	}()

	return errCh

}

func (rdbs *Rdbs) insertBatchTx(
	ctx context.Context,
	tx *sql.Tx,
	batch []map[string]interface{},
) error {

	query := rdbs.Adapter.BuildUpsertQuery(
		rdbs.Schema,
		rdbs.Upsert,
		len(batch),
	)

	values := make([]interface{}, 0, len(batch)*len(rdbs.Schema.Columns))

	for _, rec := range batch {
		for _, col := range rdbs.Schema.Columns {
			values = append(values, rec[col])
		}
	}

	_, err := tx.ExecContext(ctx, query, values...)
	return err
}


