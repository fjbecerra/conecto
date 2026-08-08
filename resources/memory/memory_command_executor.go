package memory

import (
	"conecto/core"
	"conecto/core/commands"
	"context"
	"errors"
)

type MemoryExecutor struct {
    Store *[]core.Event
}

func (m *MemoryExecutor) Execute(ctx context.Context,commands []commands.Command,
) error {

    for _, command := range commands {
        switch c := command.(type) {
            case *MemoryCommand:
                *m.Store = append(
                    *m.Store,
                    c.Events...,
                )
                return nil

            default:
                return errors.New("unsupported command")
        }
    }
    return nil    
}

func (m *MemoryExecutor) Close() error{
	return nil
}