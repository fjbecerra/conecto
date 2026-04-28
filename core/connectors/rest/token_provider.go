package rest

import "net/http"

type TokenProvider interface {
	Apply(req *http.Request) error
}

type QueryTokenProvider struct {
	ParamName string
	Token string
}

func (f *QueryTokenProvider) Apply(req *http.Request) error {
	q := req.URL.Query()
	q.Set(f.ParamName, f.Token)
	req.URL.RawQuery = q.Encode()
	return nil
}

type BearerTokenProvider struct {
	Token string
}

func (b *BearerTokenProvider) Apply(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+b.Token)
	return nil
}