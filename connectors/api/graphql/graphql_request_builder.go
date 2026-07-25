package graphql

import (
	"bytes"
	"conecto/connectors/api"
	"context"
	"encoding/json"
	"net/http"
)

type GraphQLRequestBuilder struct {
	Endpoint string
	Query    string

	VariableCursorKey string
	Headers           map[string]string
}

type graphQLBody struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

func (b *GraphQLRequestBuilder) Build(ctx context.Context, cursor *api.PageCursor) (*http.Request, error) {

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
		b.Endpoint,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}
