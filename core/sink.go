package core

import (
	"context"
)

type Sink [T any] interface {
    Write(ctx context.Context, in <-chan T) <- chan error
}