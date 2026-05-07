package memory

import (
	"conecto/core"
	"context"
)

type MemoryCommiter struct{
	stateStore MemoryStateStore
	pipelineID string
}

func (m *MemoryCommiter) Commit(ctx context.Context, batch []core.Event, state core.State) error {
	error := m.stateStore.Save(ctx, m.pipelineID, state)
	if error !=nil {
		return error
	}
	return nil
}


