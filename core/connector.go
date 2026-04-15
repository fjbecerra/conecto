package core

import (
	"context"
)

type Connector[OUTPUT any] interface {
	Run(ctx context.Context) (<-chan OUTPUT, <-chan error)
}






