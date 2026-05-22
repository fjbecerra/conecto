package memory

import (
	"conecto/core"
    "conecto/core/commands"
	"context"
)

type MemoryExecutor struct {
    Store *[]core.Event
}

func (m *MemoryExecutor) Execute(ctx context.Context,command commands.Command,
) error {

    switch c := command.(type) {

    case *MemoryCommand:

        *m.Store = append(
            *m.Store,
            c.Events...,
        )

        return nil

    default:
        panic("unsupported memory command")
    }
}

func (m *MemoryExecutor) Close() error{
	return nil
}