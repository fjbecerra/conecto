package _http

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"context"
	"net/http"
)

type Client struct{
	Client 	IClient
	Provider Provider
	CredentialService credentials.CredentialService
}

func NewClient(client IClient, provider Provider, credentialService credentials.CredentialService) *Client{
    return &Client{
        Client: client,
		Provider: provider,
		CredentialService: credentialService,
    }
}

func (c *Client) Fetch(context context.Context, req *http.Request, connection connections.Connection) (*HttpResponse, error){

	credential, err := c.CredentialService.Get(context, connection.ID)
	if err != nil {
		return nil, err
	}

	if c.Provider != nil {
		if err := c.Provider.Apply(req, credential); err != nil {
			return nil, err
		}
	}
	return c.Client.Fetch(req)
}

func (c *Client) Close() error{	
	err := c.CredentialService.Close()
	if err!=nil{
		return err
	}
    return  c.Client.Close()    
}