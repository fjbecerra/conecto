package sinks

import (
	"conecto/core"
	"conecto/core/sinks/rdbs"
	"context"
	"fmt"
	"sync"
)

type SinkMemory struct {
    mu   sync.Mutex
    Mstore *[]map[string]interface{}
    Adapter rdbs.Adapter
}

func NewMemorySink(mstore*[]map[string]interface{}, adapter rdbs.Adapter) *SinkMemory{
    var mu sync.Mutex
    return &SinkMemory{
        mu : mu,
        Mstore: mstore,
        Adapter: adapter,
    }
}

func (m *SinkMemory) WriteBatch(ctx context.Context, batch [] core.Event) error{
    fmt.Println("SINK: received event")
    for _, event := range batch {
        m.mu.Lock()
        record, error := m.Adapter.Decode(event.Payload)
        if error != nil {
            return error
        }
        *m.Mstore = append(*m.Mstore, record)
        m.mu.Unlock()
    }
    fmt.Println("SINK: exiting Run()")
    return nil
}

func (m *SinkMemory) Commit(
	ctx context.Context,
	cursor core.Cursor,
) error {

	// checkpoint store (future DB table)
	fmt.Println("checkpoint:", cursor)
	return nil
}

func (m *SinkMemory) Open(ctx context.Context) error {
	fmt.Println("SINK: open")
	return nil
}

func (m *SinkMemory) Close() error {
	fmt.Println("SINK: close")
	return nil
}


