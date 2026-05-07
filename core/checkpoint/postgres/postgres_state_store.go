package postgres

import (
	"conecto/core"
	"context"
	"database/sql"
	"encoding/json"
)

type PostgresStateStore struct {
	DB *sql.DB
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