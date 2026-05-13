package statestores

import (
	"conecto/core"
	"conecto/core/commands"
)

type StateStore interface {
	Load(runtime core.Runtime) (core.State, error)
	Save(runtime core.Runtime, state core.State) ([]commands.Command, error)
	Close()error
}

