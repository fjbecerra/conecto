package core

import (
	"context"
)

type BaseExtractor struct {
	Schema SchemaConfig
}

type Extractor [T any]interface {
    Extract(ctx context.Context, in <-chan T) (<-chan Record, <-chan error)
}

