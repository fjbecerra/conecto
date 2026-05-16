package rest

import (
	"conecto/core"
	"conecto/core/connectors/rest/auths"
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

func (c *RestClient) Fetch(runtime core.Runtime, url string) ([]byte, error) {

	token, err := c.TokenStore.Get(runtime)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(runtime.Context, "GET", url, nil)
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