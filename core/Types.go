package core

type Event struct {
	Payload []byte
	Cursor  Cursor
}

type Batch struct {
	Events []Event
	Cursor Cursor
}


type SourceType string
const (
	SourceRest 		 SourceType = "rest"
	SourceMockedRest SourceType = "mocked_rest"
)

type TransformerType string
const (
	TransformerExtractor TransformerType = "extractor"
)

type SinkType string
const (
	Rdbs  SinkType = "rdbs"
)

type RdbsType string
const (
	Postgres RdbsType = "postgres"
)
