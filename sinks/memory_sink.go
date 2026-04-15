package sinks

import (
	"conecto/core"
	"context"
	"sync"
)

type MemorySink struct {
    mu   sync.Mutex
    data []core.Record
}

func NewMemorySink() *MemorySink  {
    var mu sync.Mutex
    return &MemorySink{
        mu : mu,
        data: make([]core.Record, 0),
    }
}

func (m *MemorySink) Write(ctx context.Context, in <-chan core.Record) <- chan error {
    errCh := make(chan error, 1)

    go func() {
        defer close(errCh)

        for record := range in {
            m.mu.Lock()
            m.data = append(m.data, record)
            m.mu.Unlock()
        }
    }()

    return errCh
}

func (m *MemorySink) Data() []core.Record {
    m.mu.Lock()
    defer m.mu.Unlock()

    out := make([]core.Record, len(m.data))
    copy(out, m.data)
    return out
}