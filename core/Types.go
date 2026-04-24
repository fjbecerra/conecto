package core

type Record map[string]interface{}


type SourceType string
const (
	SourceRest 		 SourceType = "rest"
	SourceMockedRest SourceType = "mocked_rest"
)

type TransformType string
const (
	TransformExtractor TransformType = "extractor"
)

type SinkType string
const (
	Rdbs  SinkType = "rdbs"
)

type RdbsType string
const (
	Postgres RdbsType = "postgres"
)
