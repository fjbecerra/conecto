package rest

import (
	"context"
	"encoding/json"
	"net/url"
	"github.com/tidwall/gjson"
)

type PaginationProvider struct {
	Client 	IClient
	BaseUrl string
	DataPath string
	ResponseNextPath string
	RequestParam string
}

func (p *PaginationProvider) FetchPage(ctx context.Context, cursor *Cursor) (Page[json.RawMessage], error) {

	u, _ := url.Parse(p.BaseUrl)
	q := u.Query()

	if cursor != nil && p.RequestParam != "" {
		q.Set(p.RequestParam, cursor.Value)
	}

	u.RawQuery = q.Encode()

	body, err := p.Client.Fetch(ctx, u.String())
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	// extract next cursor
	next := gjson.GetBytes(body, p.ResponseNextPath).String()

	var nextCursor *Cursor
	if next != "" {
		nextCursor = &Cursor{Value: next}
	}

	res := gjson.GetBytes(body, p.DataPath)
	
	var rows []json.RawMessage
	for _, item := range res.Array() {
		rows = append(rows, json.RawMessage(item.Raw))
	}

	return Page[json.RawMessage]{
		Data:       rows,
		NextCursor: nextCursor,
		HasMore:    next != "",
	}, nil
}