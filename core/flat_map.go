package core

import (
	"conecto/core/sources"
	"context"
)

type FlatMapSource[INPUT any, OUTPUT any] struct {
	Upstream sources.Source[INPUT]
	MapFn    func(INPUT) []OUTPUT
}

func (f *FlatMapSource[INPUT, OUTPUT]) Fetch(ctx context.Context) (<-chan OUTPUT, <-chan error) {
	in, errCh := f.Upstream.Fetch(ctx)

	out := make(chan OUTPUT)

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

