package graphql

import (
	"conecto/connectors/_http"
	"github.com/tidwall/gjson"
)

type GraphQLCursorExtractor struct {
	HasMorePath string
	CursorPath string
}

func (e *GraphQLCursorExtractor) Extract(resp *_http.HttpResponse,) (*_http.PageCursor, bool) {

	hasMore := gjson.GetBytes(
		resp.Body,
		e.HasMorePath,
	).Bool()

	if !hasMore {
		return nil, false
	}

	cursor := gjson.GetBytes(
		resp.Body,
		e.CursorPath,
	).String()

	return &_http.PageCursor{
		Value: cursor,
	}, true

}