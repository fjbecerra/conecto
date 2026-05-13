package engines

import (
	"conecto/core"
	"conecto/core/commands"
	"conecto/core/retry"
	"conecto/core/sinks"
	"conecto/core/statestores"
	"context"
	"errors"
)

type AtLeastOnceCommitStrategy struct{
    SinkRetry retry.Executor
	Sink sinks.Sink
    StateStoreRetry retry.Executor
    StateStore statestores.StateStore
    Executor   commands.CommandExecutor
}

func (a *AtLeastOnceCommitStrategy) Commit(runtime core.Runtime,batch core.Batch) error {

    // 1. WRITE EVENTS

    if err := a.SinkRetry.Do(runtime.Context, func() error {

        cmds, err :=
            a.Sink.WriteBatch(runtime, batch.Events)
        if err != nil {
            return err
        }

        for _, cmd := range cmds {

            if err := a.Executor.Execute(
                runtime.Context,
                cmd,
            ); err != nil {
                return err
            }
        }

        return nil

    }); err != nil {
        return err
    }

    // 2. SAVE CHECKPOINT

    if err := a.StateStoreRetry.Do(runtime.Context, func() error {

        var status core.Status
        if batch.IsLast {
            status = core.Completed
        }
        cmds, err :=
            a.StateStore.Save(
                runtime,
                core.State{
                    Cursor: batch.Cursor,
                    Status: status,
                },
            )

        if err != nil {
            return err
        }

        for _, cmd := range cmds {

            if err := a.Executor.Execute(
                runtime.Context,
                cmd,
            ); err != nil {
                return err
            }
        }

        return nil

    }); err != nil {
        return err
    }

    return nil
}

func (e *AtLeastOnceCommitStrategy) Shutdown(
    ctx context.Context,
) error {

    var errs []error

    if err := e.StateStore.Close(); err != nil {
        errs = append(errs, err)
    }

    if err := e.Executor.Close(); err != nil {
        errs = append(errs, err)
    }

    return errors.Join(errs...)
}