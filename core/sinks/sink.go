package sinks

import (
	"conecto/core"
	"conecto/core/commands"
)


type Sink interface {
	WriteBatch(runtime core.Runtime, batch []core.Event) ([]commands.Command, error)
}