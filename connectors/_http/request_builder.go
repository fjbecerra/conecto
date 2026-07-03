package _http

import (
	"context"
	"net/http"
)

type RequestBuilder interface {
    Build(ctx context.Context, cursor *PageCursor,) (*http.Request, error)
}