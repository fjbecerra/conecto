package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"conecto/core/statestores"
)

type EventMeta map[string]any

type Event struct {
	Payload []byte
	Meta EventMeta
}

func NewEvent(payload []byte) Event {
	meta:= EventMeta{
		"__event_id": Generate(payload),
		"__ingested_at": time.Now(),
	}
	return Event{
		Payload: payload,
		Meta: meta,
	}
}

func Generate(payload []byte) string {
    h := sha256.New()
    h.Write(payload)
    return hex.EncodeToString(h.Sum(nil))
}

type Batch struct {
	Events []Event
	Cursor statestores.Cursor
	IsLast bool
}