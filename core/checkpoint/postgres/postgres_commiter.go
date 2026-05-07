package postgres

import (
	"conecto/core"
	"conecto/core/checkpoint"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Schema struct {
	Table   string
	Columns []string
}

type PostgresCommitter struct {
	DB         *sql.DB
	StateStore *PostgresStateStore
	Schema Schema
	PipelineID string
	
}

func (c *PostgresCommitter) Commit(
	ctx context.Context,
	batch []core.Event,
	state core.State,
) error {

	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := c.insertBatchTx(ctx, tx, batch); err != nil {
		tx.Rollback()
		return err
	}

	if err := c.saveStateTx(ctx, tx, state); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *PostgresCommitter) insertBatchTx(
	ctx context.Context,
	tx *sql.Tx,
	batch []core.Event,
) error {

	if len(batch) == 0 {
		return nil
	}

	query := buildInsertQuery(c.Schema, len(batch))

	values := make([]interface{}, 0, len(batch)*len(c.Schema.Columns))

	for _, ev := range batch {

		rec, err := evToMap(ev)
		if err != nil {
			return err
		}

		for _, col := range c.Schema.Columns {
			values = append(values, rec[col])
		}
	}

	_, err := tx.ExecContext(ctx, query, values...)
	return err
}

func buildInsertQuery(schema Schema, batchSize int) string {
	cols := strings.Join(schema.Columns, ", ")

	rows := make([]string, batchSize)

	for i := 0; i < batchSize; i++ {
		rows[i] = fmt.Sprintf("(%s)", placeholders(len(schema.Columns), i))
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id",
		schema.Table,
		cols,
		strings.Join(rows, ", "),
	)
}

func evToMap(ev core.Event) (checkpoint.Record, error) {
	var m map[string]interface{}

	if err := json.Unmarshal(ev.Payload, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func placeholders(cols int, row int) string {
	out := make([]string, cols)

	base := row * cols

	for i := 0; i < cols; i++ {
		out[i] = fmt.Sprintf("$%d", base+i+1)
	}

	return strings.Join(out, ", ")
}

func (c *PostgresCommitter) saveStateTx(
	ctx context.Context,
	tx *sql.Tx,
	state core.State,
) error {

	b, _ := json.Marshal(state.Cursor)

	_, err := tx.ExecContext(
		ctx,
		`
		INSERT INTO pipeline_state (
			pipeline_id,
			cursor,
			watermark
		)
		VALUES ($1,$2,$3)

		ON CONFLICT(pipeline_id)
		DO UPDATE SET
			cursor = EXCLUDED.cursor,
			watermark = EXCLUDED.watermark
		`,
		c.PipelineID,
		b,
		state.Watermark,
	)

	return err
}