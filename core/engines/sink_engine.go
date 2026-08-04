package engines

import (
	"conecto/core"
	"conecto/core/commands"
	"conecto/core/retry"
	"conecto/core/statestores"
	"context"
)


type SinkCommiter interface{
	Commit(ctx context.Context, ID string, streamName string, batch core.Batch) error
}

type SinkEngine struct {
    SinkRetry retry.Executor
	Sink core.Sink
    StateStoreRetry retry.Executor
    StateStore statestores.StateStore
    Executor   commands.CommandExecutor
}

func (a *SinkEngine) Commit(context context.Context, ID string, streamName string, batch core.Batch) error {

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
                streamName,
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
