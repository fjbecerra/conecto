package core

import (
	"conecto/core/sources"
	"context"
)

type FlatMapSource[A any, B any] struct {
	Upstream sources.Source[A]
	MapFn    func(A) []B
}

func (f *FlatMapSource[A, B]) Fetch(ctx context.Context) (<-chan B, <-chan error) {
	in, errCh := f.Upstream.Fetch(ctx)

	out := make(chan B)

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

