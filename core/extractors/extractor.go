package extractors

import (
	"context"
	"conecto/core"
)

type BaseExtractor struct {
	Config FieldsConfig
}

type Extractor [T any]interface {
    Extract(ctx context.Context, in <-chan T) (<-chan core.Record, <-chan error)
}

