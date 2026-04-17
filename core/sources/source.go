package sources

import (
	"context"
)

type Source[INPUT any] interface {
	Fetch(ctx context.Context) (<-chan INPUT, <-chan error)
}

