package rest

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/tidwall/gjson"
)

type PaginationProvider struct {
	RestClient 	RestClient
	BaseUrl string
	DataPath string
	ResponseNextPath string
	RequestParam string
}

func (p *PaginationProvider) FetchPage(context context.Context, cursor *PageCursor, ID string) (Page[json.RawMessage], error) {

	var u, _ = url.Parse(p.BaseUrl)
	q := u.Query()

	if cursor != nil && p.RequestParam != "" {
		q.Set(p.RequestParam, cursor.Value)
	}
	
	u.RawQuery = q.Encode()

	body, err := p.RestClient.Fetch(context, u.String(), ID)
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	// extract next cursor
	next := gjson.GetBytes(body, p.ResponseNextPath).String()

	var nextCursor *PageCursor
	if next != "" {
		nextCursor = &PageCursor{Value: next}
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

func (p *PaginationProvider) Close() error {
	return p.RestClient.Close()
}