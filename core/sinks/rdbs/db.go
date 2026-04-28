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
}


func (rdbs *Rdbs) WriteBatch(ctx context.Context, batch [] core.Event) error {
	fmt.Println("SINK: writing batch size =", len(batch))

	tx, err := rdbs.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rows := make([]map[string]interface{}, 0, len(batch))

	for _, ev := range batch {

		rec, err := rdbs.Adapter.Decode(ev.Payload)
		if err != nil {
			tx.Rollback()
			return err
		}

		rows = append(rows, rec)
	}

	if err := rdbs.insertBatchTx(ctx, tx, rows); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *Rdbs) Commit(
	ctx context.Context,
	cursor core.Cursor,
) error {

	// checkpoint store (future DB table)
	fmt.Println("checkpoint:", cursor)
	return nil
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

func (r *Rdbs) Open(ctx context.Context) error {
	fmt.Println("SINK: open")
	return nil
}

func (r *Rdbs) Close() error {
	fmt.Println("SINK: close")
	return nil
}


