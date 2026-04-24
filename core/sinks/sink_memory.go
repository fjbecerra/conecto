package sinks

import (
	"conecto/core"
	"context"
	"sync"
)

type SinkMemory struct {
    mu   sync.Mutex
    Mstore *[]core.Record
}

func NewMemorySink(mstore *[]core.Record) *SinkMemory{
    var mu sync.Mutex
    return &SinkMemory{
        mu : mu,
        Mstore: mstore,
    }
}

func (m *SinkMemory) Write(ctx context.Context, in <-chan core.Record) <- chan error {
    errCh := make(chan error, 1)

    go func() {
        defer close(errCh)

        for record := range in {
            m.mu.Lock()
            *m.Mstore = append(*m.Mstore, record)
            m.mu.Unlock()
        }
    }()

    return errCh
}