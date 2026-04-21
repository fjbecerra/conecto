package core

import (
	"conecto/core/sources"
	"context"
)

type FlatMapSource[IN any, OUT any] struct {
	Upstream sources.Source[IN]
	MapFn    func(IN) []OUT
}

func (f *FlatMapSource[IN, OUT]) Fetch(ctx context.Context) (<-chan OUT, <-chan error) {
	in, errCh := f.Upstream.Fetch(ctx)

	out := make(chan OUT)

	go func() {
		defer close(out)

		for v := range in {
			results := f.MapFn(v)
			for _, r := range results {
				out <- r
			}
		}
	}()

	return out, errCh
}

