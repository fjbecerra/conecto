package memory

import (
	"conecto/core"
)

type MemoryCommand struct {
    Events []core.Event
}

func (m *MemoryCommand) Kind() string {
    return "memory"
}

