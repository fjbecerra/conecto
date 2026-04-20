package transforms

import (
	"conecto/core/sources"
	"context"
	"sync"
)

type MapSource[INPUT any, OUTPUT any] struct {
	Upstream sources.Source[INPUT]
	MapFn    func(INPUT) OUTPUT
	Workers  int
}

func (m *MapSource[INPUT, OUTPUT]) Fetch(ctx context.Context) (<-chan OUTPUT, <-chan error) {
	in, errCh := m.Upstream.Fetch(ctx)

	out := make(chan OUTPUT)
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