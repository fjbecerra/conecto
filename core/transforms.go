package core

import (
	"conecto/core/sources"
	"context"
)

type MapSource[INPUT any, OUTPUT any] struct {
	Upstream sources.Source[INPUT]
	MapFn    func(INPUT) OUTPUT
}

func (m *MapSource[INPUT, OUTPUT]) Fetch(ctx context.Context) (<-chan OUTPUT, <-chan error) {
	in, errCh := m.Upstream.Fetch(ctx)

	out := make(chan OUTPUT)

	go func() {
		defer close(out)

		for v := range in {
			out <- m.MapFn(v)
		}
	}()

	return out, errCh
}