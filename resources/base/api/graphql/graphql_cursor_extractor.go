package graphql

import (
	"conecto/resources/base/api"
	"github.com/tidwall/gjson"
)

type GraphQLCursorExtractor struct {
	HasMorePath string
	CursorPath  string
}

func (e *GraphQLCursorExtractor) Extract(resp *api.HttpResponse) (*api.PageCursor, bool) {

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

	return &api.PageCursor{
		Value: cursor,
	}, true

}
