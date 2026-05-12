package statestores

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sinks/commands"
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
		`SELECT cursor FROM pipeline_state WHERE pipeline_id=$1`,
		id,
	).Scan(&cursorBytes)

	if err == sql.ErrNoRows {
		return core.State{}, nil
	}
	if err != nil {
		return core.State{}, err
	}

	json.Unmarshal(cursorBytes, &st.Cursor)
	return st, nil
}

func (c *PostgresStateStore) Save(ctx context.Context, pipelineID string,state core.State) (sinks.Command, error) {

	b, error := json.Marshal(state.Cursor)
	if error != nil{
		return nil, error
	}

	return commands.New(
		`
		INSERT INTO pipeline_state (
			pipeline_id,
			cursor
		)
		VALUES ($1,$2)

		ON CONFLICT(pipeline_id)
		DO UPDATE SET
			cursor = EXCLUDED.cursor
		`,
		pipelineID,
		b,
	),nil
}
