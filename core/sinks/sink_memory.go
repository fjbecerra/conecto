package sinks

import (
	"context"
	"sync"
)

type SinkMemory[IN any] struct {
    mu   sync.Mutex
    data []IN
}

func NewMemorySink[IN any]() *SinkMemory[IN]{
    var mu sync.Mutex
    return &SinkMemory[IN]{
        mu : mu,
        data: make([]IN, 0),
    }
}

func (m *SinkMemory[IN]) Write(ctx context.Context, in <-chan IN) <- chan error {
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

func (m *SinkMemory[IN]) Data() []IN {
    m.mu.Lock()
    defer m.mu.Unlock()

    out := make([]IN, len(m.data))
    copy(out, m.data)
    return out
}