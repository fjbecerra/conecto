package core

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
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
	Cursor Cursor
	IsLast bool
}

type Status int

const (
    Running Status = iota
    Completed
    Failed
	Stopped
)

func (s Status) String() string {
    switch s {
    case Running:
        return "RUNNING"
    case Completed:
        return "COMPLETED"
	case Failed:
        return "FAILED"
	case Stopped:
        return "STOPPED"
    default:
        return "UNKNOWN"
    }
}

func ParseStatus(s string) (Status, error) {
    switch s {
    case "RUNNING":
        return Running, nil
    case "COMPLETED":
        return Completed, nil
	case "FAILED":
        return Failed, nil
	case "STOPPED":
        return Stopped, nil
    default:
        return 0, fmt.Errorf("unknown status: %s", s)
    }
}

type State struct {
	Cursor   Cursor
	Status   Status
}

type Cursor map[string]string

func Encode(c Cursor) string {
	if c == nil {
		return ""
	}

	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) Cursor {
	if s == "" {
		return nil
	}

	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil
	}

	var c Cursor
	_ = json.Unmarshal(b, &c)

	return c
}

type Runtime struct {
	PipelineId 		string
	Provider       string
	Context 		context.Context
}
