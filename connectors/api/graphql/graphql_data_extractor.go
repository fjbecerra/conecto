package graphql

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type GraphQLDataExtractor struct {
	Path string
}

func (e *GraphQLDataExtractor) Extract(
	body []byte,
) ([]json.RawMessage, error) {

	edges := gjson.GetBytes(body, e.Path)

	rows := make([]json.RawMessage, 0, len(edges.Array()))

	for _, edge := range edges.Array() {

		node := edge.Get("node")

		rows = append(rows,
			json.RawMessage(node.Raw),
		)

	}

	return rows, nil

}