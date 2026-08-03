package api

import (
	"conecto/auth/connections"
	"context"
	"encoding/json"
)

type PaginationProvider struct {
	Builder RequestBuilder
	Client  *Client
	Data    DataExtractor
	Cursor  CursorExtractor
}

func (p *PaginationProvider) FetchPage(context context.Context, cursor *PageCursor, connection connections.Connection) (Page[json.RawMessage], error) {

	req, err := p.Builder.Build(context, cursor, connection)
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	resp, err := p.Client.Fetch(context, req, connection)
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	rows, err := p.Data.Extract(resp.Body)
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	nextCursor, hasMore := p.Cursor.Extract(resp)

	return Page[json.RawMessage]{
		Data:       rows,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (p *PaginationProvider) Close() error {
	return p.Client.Close()
}
