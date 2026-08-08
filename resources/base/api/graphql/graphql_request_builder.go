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
	//Headers           map[string]string
}

type graphQLBody struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

func (b *GraphQLRequestBuilder) Build(ctx context.Context, cursor *api.PageCursor, connection core.Connection) (*http.Request, error) {

	vars := map[string]any{}

	if cursor != nil && b.VariableCursorKey != "" {
		vars[b.VariableCursorKey] = cursor.Value
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

	// for k, v := range b.Headers {
	// 	req.Header.Set(k, v)
	// }

	return req, nil
}
