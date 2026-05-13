package commands

import (
	"conecto/core"
	"context"
)

type MemoryCommandExecutor struct {
    Store *[]core.Event
}

func (m *MemoryCommandExecutor) Execute(ctx context.Context,command Command,
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

func (m *MemoryCommandExecutor) Close() error{
	return nil
}