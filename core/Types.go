package core

type Event struct {
	Payload []byte
}

type Batch struct {
	Events []Event
	Cursor Cursor
}

type State struct {
	Cursor    Cursor
}

const(
	PipelineId string = "pipeline_id"
)