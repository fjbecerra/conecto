package core

import (
	"context"
)

type Sink interface {
    Write(ctx context.Context, in <-chan Record) <- chan error
}