package statestores

import (
	"conecto/core"
	"conecto/core/commands"
	"database/sql"
	"encoding/json"
)


type PostgresStateStore struct {
	DB *sql.DB
}

func (s *PostgresStateStore) Load(runtime core.Runtime) (core.State, error) {
	var st core.State
	var cursorBytes []byte
	var status string

	err := s.DB.QueryRowContext(runtime.Context,
		`SELECT cursor, status FROM pipeline_state WHERE pipeline_id=$1`,
		runtime.PipelineId,
	).Scan(&cursorBytes, &status)

	if err == sql.ErrNoRows {
		return core.State{}, nil
	}
	if err != nil {
		return core.State{}, err
	}

	json.Unmarshal(cursorBytes, &st.Cursor)
	stat, error := core.ParseStatus(status)
	if error!=nil {
		return core.State{}, error
	}
	st.Status=stat
	return st, nil
}

func (c *PostgresStateStore) Save(runtime core.Runtime,state core.State) ([]commands.Command, error) {

	b, error := json.Marshal(state.Cursor)
	if error != nil{
		return nil, error
	}
	query:= `
		INSERT INTO pipeline_state (
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
	values = append(values, runtime.PipelineId, b, state.Status.String())
	return []commands.Command{
        &commands.SQLCommand{
           Query: query,
		Values: values,
	}},nil
}

func (c *PostgresStateStore) Close() error{
	return c.DB.Close()
}