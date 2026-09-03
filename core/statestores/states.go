package statestores

import (
	"encoding/base64"
	"encoding/json"
)

type StateStatus string

const (
    Running StateStatus = "running"
    Completed StateStatus = "completed"
    Failed StateStatus = "failed"
	Stopped  StateStatus= "stopped"
    Idle StateStatus = "idle"
)

type State struct {
	Cursor     Cursor
	Status     StateStatus
	Watermark *string //nil means no checkpoint yet, which is the first sync
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