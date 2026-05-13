package engines

import (
	"conecto/core"
	"conecto/core/commands"
	"conecto/core/retry"
	"conecto/core/sinks"
	"conecto/core/statestores"
	"context"
	"database/sql"
	"errors"
)

type Transactional struct {
    DB *sql.DB
	Sink sinks.Sink
	StateStore statestores.StateStore
    Retry retry.Executor
}

func (t *Transactional) Commit(runtime core.Runtime, batch core.Batch) error {
			
    return t.Retry.Do(runtime.Context, func() error {

        tx, err := t.DB.BeginTx(runtime.Context, nil)
        if err != nil {
            return err
        }

        defer tx.Rollback()

        executor := &commands.SQLExecutor{
            Tx: tx,
        }

        sinkCommands, err :=
            t.Sink.WriteBatch(runtime, batch.Events)
        if err != nil {
            return err
        }

        for _, cmd := range sinkCommands {

            if err := executor.Execute(
                runtime.Context,
                cmd,
            ); err != nil {
                return err
            }
        }

        var status core.Status
        if batch.IsLast {
            status = core.Completed
        }
        stateCommands, err :=
            t.StateStore.Save(
                runtime,
                core.State{
                    Cursor: batch.Cursor,
                    Status: status,
                },
            )

        if err != nil {
            return err
        }

        for _, cmd := range stateCommands {

            if err := executor.Execute(
                runtime.Context,
                cmd,
            ); err != nil {
                return err
            }
        }

        return tx.Commit()
    })
}

func (t *Transactional) Shutdown(
    ctx context.Context,
) error {

    var errs []error

    if err := t.StateStore.Close(); err != nil {
        errs = append(errs, err)
    }

    if err := t.DB.Close(); err != nil {
        errs = append(errs, err)
    }

    return errors.Join(errs...)
}