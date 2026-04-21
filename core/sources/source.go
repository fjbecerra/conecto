package sources

import (
	"context"
)

type Source[OUT any] interface {
	Fetch(ctx context.Context) (<-chan OUT, <-chan error)
}

