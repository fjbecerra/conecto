package shopify

import (
	"bytes"
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Connector struct {
	ClientID string
	ClientSecret string
	Scopes []string
	HttpClient *http.Client
}



func (c Connector) Name() string {
	return "shopify"
}

func (c Connector) Exchange(ctx context.Context,connection connections.Connection,code string)(credentials.Credential,error){
	shop := connection.Metadata["shop"]
	body := map[string]string{
		"client_id": c.ClientID,
		"client_secret": c.ClientSecret,
		"code": code,
	}
	// Convert to JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	response,err :=c.HttpClient.Post(
			fmt.Sprintf("https://%s/admin/oauth/access_token",shop),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
	defer response.Body.Close()


	if err != nil {
		return credentials.Credential{},err
	}

	return credentials.Credential{

		Type:"oauth2",

		Data:map[string]string{
			"access_token": response.Header.Get("X-Shopify-Access-Token"),
		},
	},nil
}

func (c Connector) AuthorizeURL(ctx context.Context, connection connections.Connection,state string)(string,error){
	shop := connection.Metadata["shop"]
	params := url.Values{}
	
	params.Set(
		"client_id",c.ClientID,
	)
	params.Set(
		"scope",strings.Join(c.Scopes,",",),
	)

	params.Set(
		"redirect_uri", "https://yourapp.com/oauth/callback",
	)

	params.Set(
		"state", state,
	)

	return fmt.Sprintf("https://%s/admin/oauth/authorize?%s", shop, params.Encode()),nil
}

func (c Connector) Refresh(ctx context.Context,connection connections.Connection,credential credentials.Credential) (credentials.Credential,error){
	return credential,nil
}