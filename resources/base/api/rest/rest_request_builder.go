package rest

import (
	"conecto/core"
	"conecto/resources/base/api"
	"context"
	"net/http"
	"net/url"
)

type RestRequestBuilder struct {
	EndpointProvider  api.EndPointProvider
	Method      string
	CursorParam string
	Headers     map[string]string
}

func (b *RestRequestBuilder) Build(ctx context.Context, cursor *api.PageCursor, connection core.Connection) (*http.Request, error) {

	u, err := url.Parse(b.EndpointProvider.Apply(connection))
	if err != nil {
		return nil, err
	}

	q := u.Query()

	// inject cursor
	if cursor != nil && b.CursorParam != "" {
		q.Set(b.CursorParam, cursor.Value)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		b.Method,
		u.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// headers
	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}
