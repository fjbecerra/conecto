package core

import (
	"conecto/core/sources"
	"context"
)

type MapSource[INPUT any, OUTPUT any] struct {
	Upstream sources.Source[INPUT]
	MapFn    func(INPUT) OUTPUT
}

func (m *MapSource[A, B]) Fetch(ctx context.Context) (<-chan B, <-chan error) {
	in, errCh := m.Upstream.Fetch(ctx)

	out := make(chan B)

	go func() {
		defer close(out)

		for v := range in {
			out <- m.MapFn(v)
		}
	}()

	return out, errCh
}