package graphql

import (
	"bytes"
	"conecto/core"
	"conecto/resources/base/api"
	"context"
	"encoding/json"
	"net/http"
)

type GraphQLRequestBuilder struct {
	EndpointProvider  api.EndPointProvider
	Query    string
	VariableCursorKey string
	WatermarkPath string
	IncremenatalSyncProvider api.IncrementalSyncProvider

	
}

type graphQLBody struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

func (b *GraphQLRequestBuilder) Build(ctx context.Context, cursor *api.PageCursor, connection core.Connection, watermark *string) (*http.Request, error) {

	vars := map[string]any{}

	if cursor != nil && b.VariableCursorKey != "" {
		vars[b.VariableCursorKey] = cursor.Value
	}

	if(watermark!=nil && b.WatermarkPath != ""){
		vars[b.WatermarkPath] = b.IncremenatalSyncProvider.Apply(watermark)
	}

	payload := graphQLBody{
		Query:     b.Query,
		Variables: vars,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		b.EndpointProvider.Apply(connection),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
