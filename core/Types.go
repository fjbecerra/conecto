package core

import "time"

type Event struct {
	Payload []byte
	Cursor  Cursor
	Timestamp time.Time //to get the watermark
}

type Batch struct {
	Events []Event
	Cursor Cursor
}

type State struct {
	Cursor    Cursor
	Watermark time.Time
}

