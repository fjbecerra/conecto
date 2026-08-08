package api

import (
	"conecto/stores/credentials"
	"errors"
	"net/http"
)

type Provider interface {
	Apply(req *http.Request, credential credentials.Credential) error
}

type QueryProvider struct {
	Param string
}

func (q *QueryProvider) Apply(req *http.Request, credential credentials.Credential) error {
	if _, ok := credential.Data[q.Param]; !ok {
		return errors.New("credential does not contain the required query parameter")
	}

	values := req.URL.Query()
	values.Set(
		q.Param,
		credential.Data[q.Param],
	)
	req.URL.RawQuery = values.Encode()
	return nil
}

type BearerProvider struct {
	Key string
}

func (b *BearerProvider) Apply(req *http.Request, credential credentials.Credential) error {
	if _, ok := credential.Data[b.Key]; !ok {
		return errors.New("credential does not contain the required bearer token")
	}
	token := credential.Data[b.Key]
	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)
	return nil
}

type HeaderProvider struct {
	Name string
}

func (h *HeaderProvider) Apply(req *http.Request, credential credentials.Credential) error {
	if _, ok := credential.Data[h.Name]; !ok {
		return errors.New("credential does not contain the required header value")
	}
	req.Header.Set(
		h.Name,
		credential.Data[h.Name],
	)
	return nil
}
