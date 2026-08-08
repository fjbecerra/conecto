package engines

import (
	"conecto/core"
	"conecto/core/commands"
	"conecto/core/retry"
	"conecto/core/statestores"
	"context"
)

type Sinker struct {
    Sink core.Sink
    Executor   commands.CommandExecutor
}

type SinkCommiter interface{
	Commit(ctx context.Context, ID string, streamName string, batch core.Batch) error
}

type SinkEngine struct {
    Sinker Sinker
    SinkRetry retry.Executor	
    StateStore statestores.StateStore
   
}

func (a *SinkEngine) Commit(context context.Context, ID string, streamName string, batch core.Batch) error {


    if err := a.SinkRetry.Do(context, func() error {
        //write events
        cmdSinker, err :=
            a.Sinker.Sink.WriteBatch(context,ID, batch.Events)
        if err != nil {
            return err
        }

        // save checkpoints
        var status statestores.Status
        if batch.IsLast {
            status = statestores.Completed
        }
        cmdStateStore, err := a.StateStore.Save(
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
        cmds := append(cmdSinker, cmdStateStore...)

        if err := a.Sinker.Executor.Execute(context,cmds,); err != nil {
                return err
        }
        return nil

    }); err != nil {
        return err
    }
    return nil
}
