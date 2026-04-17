package rest

import (
	"context"
	"encoding/json"
	"net/url"
	"github.com/tidwall/gjson"
	"conecto/core"
)



type PaginationProvider struct {
	Client 	  IClient
	Config    ResponseConfig
}

func NewPaginationProvider(client IClient, configPath string) PaginationProvider{
	
	return PaginationProvider{
		Client: client,
		Config: core.LoadConfig[ResponseConfig](configPath),
	}
}

func (p *PaginationProvider) FetchPage(ctx context.Context, cursor *Cursor) (Page[json.RawMessage], error) {

	u, _ := url.Parse(p.Config.BaseUrl)
	q := u.Query()

	if cursor != nil && p.Config.Pagination.Request.Param != "" {
		q.Set(p.Config.Pagination.Request.Param, cursor.Value)
	}

	u.RawQuery = q.Encode()

	body, err := p.Client.Fetch(ctx, u.String())
	if err != nil {
		return Page[json.RawMessage]{}, err
	}

	// extract next cursor
	next := gjson.GetBytes(body, p.Config.Pagination.Response.Next.Path).String()

	var nextCursor *Cursor
	if next != "" {
		nextCursor = &Cursor{Value: next}
	}

	res := gjson.GetBytes(body, p.Config.Data.Path)
	
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