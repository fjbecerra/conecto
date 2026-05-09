package core

import "time"

type Event struct {
	Payload []byte
}

type Batch struct {
	Events []Event
	Cursor Cursor
}

type State struct {
	Cursor    Cursor
	Watermark time.Time
}

