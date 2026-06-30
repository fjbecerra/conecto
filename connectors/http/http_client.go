package http

import (
	"io"
	"net/http"
)

type IClient interface {
	Fetch(req *http.Request) (*HttpResponse, error)
	Close() error
}

type HttpClient struct {
	Client *http.Client
}

func (c *HttpClient) Fetch(req *http.Request) (*HttpResponse, error) {	

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Body:    body,
		Headers: resp.Header,
		Status:  resp.StatusCode,
	}, nil
}

func (c *HttpClient) Close() error{
	if c.Client != nil {
        c.Client.CloseIdleConnections()
    }

    return nil
}