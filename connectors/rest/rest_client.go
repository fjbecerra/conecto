package rest

import (
	"conecto/connectors/rest/auths"
	"context"
	"net/http"
)

type RestClient struct{
	Client 	IClient
	TokenProvider auths.TokenProvider
	TokenStore auths.TokenStore
}

func NewRestClient(client IClient, tokenProvicer auths.TokenProvider, tokenStore auths.TokenStore) *RestClient{
    return &RestClient{
        Client: client,
		TokenProvider: tokenProvicer,
		TokenStore: tokenStore,
    }
}

func (c *RestClient) Fetch(context context.Context, url string, ID string) ([]byte, error) {

	token, err := c.TokenStore.Get(context, ID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context, "GET", url, nil)
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

func (c *RestClient) Close() error{	
	err := c.TokenStore.Close()
	if err!=nil{
		return err
	}
    return  c.Client.Close()    
}