package memory

import (
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/core/statestores"
	"conecto/shared/config"
)

type MemoryResource struct{
	store *[]map[string]any
	retryExecutor *retry.Executor
	stateStore statestores.StateStore
	cfg MemoryResourceConfig

}

func NewMemoryResource(
	store *[]map[string]any,
	retryExecutor *retry.Executor,
	stateStore statestores.StateStore,
) *MemoryResource{
	return &MemoryResource{
		store: store,
		retryExecutor: retryExecutor,
		stateStore: stateStore,
	}
}

func (m *MemoryResource) Close() error{
	return nil
}

func (m *MemoryResource) Connector(cfg config.ConfigBytes) engines.ConnectorRunnable {
    return nil
}

func(m*MemoryResource) Sink(cfg config.ConfigBytes, fieldSpecs config.FieldsSpecs)  engines.SinkCommiter{
	memorySink := MemorySink{
		store: *m.store,
		retryExecutor: m.retryExecutor,
		stateStore: m.stateStore,
	}
    return  CreateMemorySink(memorySink, nil)
}