package core

import (
	"conecto/core/sources"
	"context"
	"sync"
)

type Transform[IN any,OUT any] func(IN) OUT

type MapSource[IN any, OUT any] struct {
	Upstream sources.Source[IN]
	MapFn    Transform[IN,OUT]
	Workers  int
}

func (m *MapSource[IN, OUT]) Fetch(ctx context.Context) (<-chan OUT, <-chan error) {
	in, errCh := m.Upstream.Fetch(ctx)

	out := make(chan OUT)
	outErr := make(chan error, 1)
	var wg sync.WaitGroup
	for i:=0; i<m.Workers; i++{
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				result:= m.MapFn(v)
				select {
					case out <- result:
					case <-ctx.Done():
						return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
		close(outErr)
	}()

	return out, mergeErrors(errCh, outErr)
}

func mergeErrors(chs ...<-chan error) <-chan error {
	out := make(chan error, 1)

	go func() {
		defer close(out)
		for _, ch := range chs {
			for err := range ch {
				if err != nil {
					out <- err
					return
				}
			}
		}
	}()

	return out
}