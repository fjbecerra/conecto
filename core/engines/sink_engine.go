package engines

import (
	"conecto/core"
	"conecto/core/commands"
	"conecto/core/retry"
	"conecto/core/statestores"
	"context"
	"errors"
)


type SinkCommiter interface{
	Commit(ctx context.Context, ID string, batch core.Batch) error
 	Shutdown(ctx context.Context) error
}

type SinkEngine struct {
    SinkRetry retry.Executor
	Sink core.Sink
    StateStoreRetry retry.Executor
    StateStore statestores.StateStore
    Executor   commands.CommandExecutor
}

func (a *SinkEngine) Commit(context context.Context, ID string, batch core.Batch) error {

    // 1. WRITE EVENTS

    if err := a.SinkRetry.Do(context, func() error {

        cmds, err :=
            a.Sink.WriteBatch(context,ID, batch.Events)
        if err != nil {
            return err
        }

        for _, cmd := range cmds {

            if err := a.Executor.Execute(
                context,
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

    if err := a.StateStoreRetry.Do(context, func() error {

        var status statestores.Status
        if batch.IsLast {
            status = statestores.Completed
        }
        cmds, err :=
            a.StateStore.Save(
                context,
                ID,
                statestores.State{
                    Cursor: batch.Cursor,
                    Status: status,
                },
            )

        if err != nil {
            return err
        }

        for _, cmd := range cmds {

            if err := a.Executor.Execute(
                context,
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

func (e *SinkEngine) Shutdown(
    ctx context.Context,
) error {

    var errs []error

    // if err := e.StateStore.Close(); err != nil {
    //     errs = append(errs, err)
    // }

    if err := e.Executor.Close(); err != nil {
        errs = append(errs, err)
    }

    return errors.Join(errs...)
}