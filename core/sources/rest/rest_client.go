package rest

import (
	"context"
	"io"
	"net/http"
)

type IClient interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

type RestClient struct{
	client *http.Client
}

func NewRestClient(client *http.Client) *RestClient{
    return &RestClient{
        client: client,
    }
}

func (c *RestClient) Fetch(ctx context.Context, url string) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
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
