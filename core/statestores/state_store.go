package statestores

import (
	"conecto/core/commands"
	"context"
)

type StateStore interface {
	Load(runtime context.Context, ID string) (State, error)
	Save(runtime context.Context, ID string, state State) ([]commands.Command, error)
	Close()error
}

