package api

import (
	"conecto/core"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type PaginationProvider struct {
	Builder RequestBuilder
	Client  *Client
	Data    DataExtractor
	Cursor  CursorExtractor
	ResponseProvider ResponseProvider
}

func (p *PaginationProvider) FetchPage(context context.Context, cursor *PageCursor, connection core.Connection, watermark *string) (Page[json.RawMessage], error) {

	req, err := p.Builder.Build(context, cursor, connection, watermark)
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	resp, err := p.Client.Fetch(context, req, connection)
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	if resp.Status < 200 || resp.Status >= 300 {
		return Page[json.RawMessage]{}, 
		 errors.New(fmt.Sprintf("%s",resp.Body))
	}

	body, err := p.ResponseProvider.Apply(resp.Body)

	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	rows, err := p.Data.Extract(body)
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

