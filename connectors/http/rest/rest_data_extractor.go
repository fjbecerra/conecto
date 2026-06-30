package rest

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type RestDataExtractor struct {
	Path string
}

func (e *RestDataExtractor) Extract(body []byte,) ([]json.RawMessage, error) {

	res := gjson.GetBytes(body, e.Path)

	rows := make([]json.RawMessage, 0, len(res.Array()))

	for _, item := range res.Array() {
		rows = append(rows, json.RawMessage(item.Raw))
	}

	return rows, nil
}