package shopify

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/resources/base/api"
	"conecto/shared/clients"
	"conecto/shared/config"
	"conecto/stores/credentials"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ShopifyResource struct {
    name              	string
    oauth             	*ShopifyOAuth
	client		 		*api.HttpClient
    credentialService 	credentials.CredentialService
    retryExecutor 		*retry.Executor
}


func NewShopifyResource(
    client		 		*clients.HttpClient,
    credentialService 	credentials.CredentialService,
	cfg 				ShopifyResourceConfig,
	retryExecutor 		*retry.Executor,
) *ShopifyResource {
    hclient:= client.Get()
    oauth := ShopifyOAuth{
        shopifyResourceConfig: cfg,
		httpClient: hclient,
	}
    return &ShopifyResource{
        oauth:   &oauth,
        client: client.Get(),
        credentialService: credentialService,
        retryExecutor: retryExecutor,
    }
}

func(s *ShopifyResource) Close() error {
	return s.client.Close()
}

func (s *ShopifyResource) Connector(cfg config.ConfigBytes) engines.ConnectorRunnable {
    connectorConfig,_ := config.Unmarshal[ShopifyConnectorConfig](cfg, config.FormatJSON)
    shopifyConnector:= ShopifyConnector{
        httpClient : s.client,
        credentialService: s.credentialService,
        retryExecutor: s.retryExecutor,
        cfg: connectorConfig,
    }
    return CreateShopifyConnector(shopifyConnector)
}

func(s*ShopifyResource) Sink(cfg config.ConfigBytes, fieldSpecs config.FieldsSpecs)  engines.SinkCommiter{
    return nil
}

type ShopifyOAuth struct {
	shopifyResourceConfig ShopifyResourceConfig
	httpClient *api.HttpClient
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
    ExpiresIn            int    `json:"expires_in"`
    RefreshToken         string `json:"refresh_token"`
    RefreshTokenExpiresIn int   `json:"refresh_token_expires_in"`
}

func (s*ShopifyResource) Name() string {
	return "shopify"
}

func (s*ShopifyResource) Exchange(ctx context.Context,connection core.Connection,code string)(credentials.Credential,error){
	shop := connection.Metadata["shop"].(string)
	bodyReq := map[string]string{
		"client_id": s.oauth.shopifyResourceConfig.ClientID,
		"client_secret": s.oauth.shopifyResourceConfig.ClientSecret,
		"code": code,
        "expiring" : "1",
	}

	return s.fetchToken(bodyReq, shop)
}

func (s*ShopifyResource) AuthorizeURL(ctx context.Context, connection core.Connection, state string) (string,error){
	shop := connection.Metadata["shop"]
	params := url.Values{}
	
	params.Set(
		"client_id",s.oauth.shopifyResourceConfig.ClientID,
	)
	params.Set(
		"scope",strings.Join(s.oauth.shopifyResourceConfig.Scopes,",",),
	)

	params.Set(
		"redirect_uri", fmt.Sprintf("%s/oauth/callback", s.oauth.shopifyResourceConfig.AppURL),
	)

	params.Set(
		"state", state,
	)

	return fmt.Sprintf("https://%s.myshopify.com/admin/oauth/authorize?%s", shop, params.Encode()),nil
}

func (s *ShopifyResource) Refresh(ctx context.Context, connection core.Connection, credential credentials.Credential) (credentials.Credential, error) {
    shop := connection.Metadata["shop"].(string)
	bodyReq := map[string]string{
        "grant_type": "refresh_token",
		"client_id": s.oauth.shopifyResourceConfig.ClientID,
		"client_secret": s.oauth.shopifyResourceConfig.ClientSecret,
        "expiring" : "1",
	}
    refreshToken := credential.Data["refresh_token"]
    if refreshToken == "" {
        return credentials.Credential{}, errors.New("missing Shopify refresh token")
    }
    bodyReq["refresh_token"] = credential.Data["refresh_token"]
    return s.fetchToken(bodyReq, shop)
}

func (s *ShopifyResource) fetchToken(bodyReq map[string]string, shop string) (credentials.Credential, error){
    
    response, err := s.oauth.httpClient.Post(fmt.Sprintf("https://%s.myshopify.com/admin/oauth/access_token",shop), bodyReq)


	if err != nil {
		return credentials.Credential{},err
	}

	if response.Status != http.StatusOK {
		return credentials.Credential{}, error(fmt.Errorf("shopify returned %d: %s", response.Status, string(response.Body)))
	}

	var tokenResp TokenResponse

	if err := json.Unmarshal(response.Body, &tokenResp); err != nil {
		return credentials.Credential{}, err
	}

    expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return credentials.Credential{

		Type:"oauth2",

		Data:map[string]string{
			"X-Shopify-Access-Token": tokenResp.AccessToken,
            "refresh_token":          tokenResp.RefreshToken,
		},
        Expiry: &expiry,
	},nil
}
