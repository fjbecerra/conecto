package core

import (
	"conecto/core/sources"
	"context"
)


type FilterSource[T any] struct {
	Upstream sources.Source[T]
	Predicate func(T) bool
}

func (f *FilterSource[T]) Fetch(ctx context.Context) (<-chan T, <-chan error) {
	in, errCh := f.Upstream.Fetch(ctx)

	out := make(chan T)

	go func() {
		defer close(out)

		for v := range in {
			if f.Predicate(v) {
				out <- v
			}
		}
	}()

	return out, errCh
}