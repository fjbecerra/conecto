package committers

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sinks/commands"
	"conecto/core/sinks/statestores"
	"conecto/core/transformers"
	"context"
	"database/sql"
)

type DBCommitter struct{
	DB      	*sql.DB	
	Sink 		sinks.Sink
	StateStore  statestores.StateStore
}

func NewDBCommiter(db *sql.DB, sink sinks.Sink, stateStore statestores.StateStore) *DBCommitter{
	return &DBCommitter{
		DB: db,
		Sink: sink,
		StateStore: stateStore,
	}
}

func (c * DBCommitter) CommitBatch(ctx context.Context, pipelineID string, batch core.Batch, transformer transformers.Transformer) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	
	if err != nil {
		return err
	}
	defer tx.Rollback()

	events, err := transformer.Transform(ctx,batch.Events,
	)
	if err != nil {
		return err
	}

	sinkCommand, err := c.Sink.WriteBatch(ctx, events)
	if err != nil {
		return err
	}

	sinkSqlCommand:= sinkCommand.(*commands.SqlCommand)
	_, err = tx.ExecContext(ctx, sinkSqlCommand.Query, sinkSqlCommand.Values...)
	if err != nil {
		return err
	}

	state := core.State{
		Cursor: batch.Cursor,
	}
	stateStoreCommand, err := c.StateStore.Save(ctx, pipelineID, state)
	if err != nil {		
		return err
	}
	stateStoreSqlCommand := stateStoreCommand.(*commands.SqlCommand)
	_, err = tx.ExecContext(ctx, stateStoreSqlCommand.Query, stateStoreSqlCommand.Values...)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func(c *DBCommitter) close() error {
	return nil
}