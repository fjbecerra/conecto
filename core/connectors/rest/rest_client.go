package rest

import (
	"context"
	"io"
	"net/http"
)

type RestClient struct{
	client *http.Client
	tokenProvider TokenProvider
}

func NewRestClient(client *http.Client, tokenProvicer TokenProvider) *RestClient{
    return &RestClient{
        client: client,
		tokenProvider: tokenProvicer,
    }
}

func (c *RestClient) Fetch(ctx context.Context, url string) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.tokenProvider != nil {
		if err := c.tokenProvider.Apply(req); err != nil {
			return nil, err
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func (c *RestClient) Close() error{
	if c.client != nil {
        c.client.CloseIdleConnections()
    }

    return nil
}