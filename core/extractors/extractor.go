package extractors

import (
	"conecto/core"
)

type BaseExtractor struct {
	Config FieldsConfig
}

type Extractor [T any]interface {
    Extract(in T) (core.Record, error)
}

