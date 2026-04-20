package sinks

import (
	"context"
)

type Sink [T any] interface {
    Write(ctx context.Context, in <-chan T) <- chan error
}