package states

import (
	"conecto/core/commands"
	"conecto/core/statestores"
	"context"
	"database/sql"
	"encoding/json"
	"conecto/sinks/base/db"
)


type PostgresStateStore struct {
	db *sql.DB
}

func NewStateStore(db *sql.DB) *PostgresStateStore{
	return &PostgresStateStore{
		db: db,
	}
}

func (s *PostgresStateStore) Load(context context.Context, ID string) (statestores.State, error) {
	var st statestores.State
	var cursorBytes []byte
	var status string

	err := s.db.QueryRowContext(context,
		`SELECT cursor, status FROM streams_state WHERE pipeline_id=$1`,
		ID,
	).Scan(&cursorBytes, &status)

	if err == sql.ErrNoRows {
		return statestores.State{}, nil
	}
	if err != nil {
		return statestores.State{}, err
	}

	json.Unmarshal(cursorBytes, &st.Cursor)
	stat, error := statestores.ParseStatus(status)
	if error!=nil {
		return statestores.State{}, error
	}
	st.Status=stat
	return st, nil
}

func (c *PostgresStateStore) Save(context context.Context, ID string,state statestores.State) ([]commands.Command, error) {

	b, error := json.Marshal(state.Cursor)
	if error != nil{
		return nil, error
	}
	query:= `
		INSERT INTO streams_state (
			pipeline_id,
			cursor,
			status
		)
		VALUES ($1,$2,$3)

		ON CONFLICT(pipeline_id)
		DO UPDATE SET
			cursor = EXCLUDED.cursor,
			status = EXCLUDED.status,
			updated_at = NOW();
		`
	values := []interface{}{}
	values = append(values, ID, b, state.Status.String())
	return []commands.Command{
        &db.SQLCommand{
           Query: query,
		Values: values,
	}},nil
}

func (c *PostgresStateStore) Close() error{
	return c.db.Close()
}