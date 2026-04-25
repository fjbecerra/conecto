package sinks

import (
	"context"
)

type Sink [IN any] interface {
    Write(ctx context.Context, in <-chan IN) (<-chan error, <-chan struct{})
}