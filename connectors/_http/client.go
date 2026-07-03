package _http

import (
	"conecto/connectors/_http/auths"
	"context"
	"net/http"
)

type Client struct{
	Client 	IClient
	TokenProvider auths.TokenProvider
	TokenStore auths.TokenStore
}

func NewClient(client IClient, tokenProvicer auths.TokenProvider, tokenStore auths.TokenStore) *Client{
    return &Client{
        Client: client,
		TokenProvider: tokenProvicer,
		TokenStore: tokenStore,
    }
}

func (c *Client) Fetch(context context.Context, req *http.Request, ID string) (*HttpResponse, error){

	token, err := c.TokenStore.Get(context, ID)
	if err != nil {
		return nil, err
	}

	if c.TokenProvider != nil {
		if err := c.TokenProvider.Apply(req, token); err != nil {
			return nil, err
		}
	}
	return c.Client.Fetch(req)
}

func (c *Client) Close() error{	
	err := c.TokenStore.Close()
	if err!=nil{
		return err
	}
    return  c.Client.Close()    
}