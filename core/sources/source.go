package sources

import (
	"context"
)

type Source[T any] interface {
	Fetch(ctx context.Context) (<-chan T, <-chan error)
}

