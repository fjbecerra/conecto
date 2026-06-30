package auths

import "net/http"

type TokenProvider interface {
	Apply(req *http.Request, token Token) error
}

type QueryTokenProvider struct {
	ParamName string
}

func (f *QueryTokenProvider) Apply(req *http.Request, token Token) error {
	q := req.URL.Query()
	q.Set(f.ParamName, token.AccessToken)
	req.URL.RawQuery = q.Encode()
	return nil
}

type BearerTokenProvider struct {}

func (b *BearerTokenProvider) Apply(req *http.Request, token Token) error {
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	return nil
}

type HeaderTokenProvider struct {
	HeaderName string
}

func (h *HeaderTokenProvider) Apply(req *http.Request, token Token) error {
	req.Header.Set(h.HeaderName, token.AccessToken)
	return nil
}