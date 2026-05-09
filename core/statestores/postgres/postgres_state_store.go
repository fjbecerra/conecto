package postgres

import (
	"conecto/core"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Record map[string]interface{}


type Schema struct {
	Table   string
	Columns []string
}

type PostgresStateStore struct {
	DB *sql.DB
	Schema Schema
}

func (s *PostgresStateStore) Load(ctx context.Context, id string) (core.State, error) {
	var st core.State
	var cursorBytes []byte

	err := s.DB.QueryRowContext(ctx,
		`SELECT cursor, watermark FROM pipeline_state WHERE pipeline_id=$1`,
		id,
	).Scan(&cursorBytes, &st.Watermark)

	if err == sql.ErrNoRows {
		return core.State{}, nil
	}
	if err != nil {
		return core.State{}, err
	}

	json.Unmarshal(cursorBytes, &st.Cursor)
	return st, nil
}

func (s *PostgresStateStore) Save(
	ctx context.Context,
	pipelineID string,
	st core.State,
) error {

	return nil
}

func (c *PostgresStateStore) Commit(
	ctx context.Context,
	pipelineID string,
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

	if err := c.saveStateTx(ctx, pipelineID, tx, state); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *PostgresStateStore) insertBatchTx(
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

func evToMap(ev core.Event) (Record, error) {
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

func (c *PostgresStateStore) saveStateTx(
	ctx context.Context,
	pipelineID string,
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
		pipelineID,
		b,
		state.Watermark,
	)

	return err
}
