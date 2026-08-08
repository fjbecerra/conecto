package api

import (
	"conecto/core"
	"conecto/stores/credentials"
	"context"
	"net/http"
)

type Client struct {
	Client            IClient
	Provider          Provider
	CredentialService credentials.CredentialService
	CredentialRefresher credentials.CredentialRefresher

}

func NewClient(client IClient, provider Provider, credentialService credentials.CredentialService) *Client {
	return &Client{
		Client:            client,
		Provider:          provider,
		CredentialService: credentialService,
	}
}

func (c *Client) Fetch(context context.Context, req *http.Request, connection core.Connection) (*HttpResponse, error) {

	credential, err := c.CredentialService.Get(context, connection)
	if err != nil {
		return nil, err
	}

	if credential.IsExpired() {
		credential, err = c.CredentialService.GetValid(context,connection,c.CredentialRefresher)
	}

	if c.Provider != nil {
		if err := c.Provider.Apply(req, credential); err != nil {
			return nil, err
		}
	}
	return c.Client.Fetch(req)
}

