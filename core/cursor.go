package core

import (
	"encoding/base64"
	"encoding/json"
)

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