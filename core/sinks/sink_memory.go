package sinks

import (
	"conecto/core"
	"conecto/core/sinks/codecs"
	"context"
	"fmt"
	"sync"
)

type SinkMemory struct {
    mu   sync.Mutex
    Mstore []map[string]interface{}
    Codec codecs.Codec
    OnWrite func()
}

func NewMemorySink(mstore[]map[string]interface{}, codec codecs.Codec) *SinkMemory{
    var mu sync.Mutex
    return &SinkMemory{
        mu : mu,
        Mstore: mstore,
        Codec: codec,
    }
}

func (m *SinkMemory) WriteBatch(ctx context.Context, batch [] core.Event) error{
    fmt.Println("SINK: received event")
     if m.OnWrite != nil {
            m.OnWrite()
        }
    for _, event := range batch {
        m.mu.Lock()
        record, error := m.Codec.Decode(event.Payload)
        
        if error != nil {
            return error
        }
        
        m.Mstore = append(m.Mstore, record)
       
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


