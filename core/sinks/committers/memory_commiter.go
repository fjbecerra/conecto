package committers

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/core/sinks/statestores"
	"conecto/core/transformers"
	"context"
)

type MemoryCommitter struct{
	Sink 		sinks.Sink
	StateStore  statestores.StateStore
}

func NewMemoryCommiter(sink sinks.Sink, stateStore statestores.StateStore) *MemoryCommitter{
	return &MemoryCommitter{
		Sink: sink,
		StateStore: stateStore,
	}
}

func (m * MemoryCommitter) CommitBatch(ctx context.Context, pipelineID string, batch core.Batch, transformer transformers.Transformer) error {
	
	_, err := m.Sink.WriteBatch(ctx, batch.Events)
	if err!=nil{
		return err
	}
	state := core.State{
		Cursor: batch.Cursor,
	}
	_, err = m.StateStore.Save(ctx, pipelineID, state)
	return err
}

func(m *MemoryCommitter) close() error {
	return nil
}