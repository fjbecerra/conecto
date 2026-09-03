package shopify

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/resources/base/api"
	"conecto/resources/base/api/graphql"
	"conecto/stores/credentials"
	"encoding/json"
	"fmt"
)


type ShopifyConnector struct {
    httpClient        	*api.HttpClient
    credentialService 	credentials.CredentialService
	retryExecutor 		*retry.Executor
	cfg 				ShopifyConnectorConfig
}

func CreateShopifyConnector(shopifyConnector ShopifyConnector) engines.ConnectorRunnable {
	name:= shopifyConnector.cfg.Name
	query:=shopifyConnector.cfg.Query

	provider:= &api.HeaderProvider{
		Name: "X-Shopify-Access-Token",
	}
	
	builder:= &graphql.GraphQLRequestBuilder{
		EndpointProvider:  &ShopifyEndpointProvider{},
		Query:    query,
		VariableCursorKey: "after",
		WatermarkPath: "query",
		IncremenatalSyncProvider: &ShopifyIncrementalSyncProvider{},

	}

	dataExtractor := &graphql.GraphQLDataExtractor{
		Path: fmt.Sprintf("data.%s.edges", name),
	}

	cursorExtractor := &graphql.GraphQLCursorExtractor{
		HasMorePath: fmt.Sprintf("data.%s.pageInfo.hasNextPage", name),
		CursorPath:  fmt.Sprintf("data.%s.pageInfo.endCursor", name),
	}

	client := *api.NewClient(shopifyConnector.httpClient, provider, shopifyConnector.credentialService)

	paginationProvider := api.PaginationProvider{
		Client:  &client,
		Builder: builder,
		Data:    dataExtractor,
		Cursor:  cursorExtractor,
		ResponseProvider: &ShopifyResponseProvider{},
	}
	connector:= &api.HttpConnector{
		Provider: &paginationProvider,
	}	
	
	return &engines.ConnectorEngine{
		Connector: connector,
		Retry:     *shopifyConnector.retryExecutor,
	}

}

type ShopifyError struct {
	Message string `json:"message"`
}

type ShopifyResponse struct {
	Errors []ShopifyError `json:"errors"`
}

type ShopifyResponseProvider struct {}

func (p *ShopifyResponseProvider) Apply(body []byte) ([]byte,error){
	var shopifyResp ShopifyResponse

	err := json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, err
	}

	if len(shopifyResp.Errors) > 0 {
		return nil, fmt.Errorf(
			"shopify error: %s",
			shopifyResp.Errors[0].Message,
		)
	}
	return body, nil
}

type ShopifyEndpointProvider struct {}

func (p *ShopifyEndpointProvider) Apply(connection core.Connection) string {
	shop := connection.Metadata["shop"]
	apiVersion := connection.Metadata["api_version"]
	return fmt.Sprintf("https://%s.myshopify.com/admin/api/%s/graphql.json", shop, apiVersion)
}

type ShopifyIncrementalSyncProvider struct {}
func (s *ShopifyIncrementalSyncProvider) Apply(watermark *string) string {
	return fmt.Sprintf("updated_at>%s", *watermark)
}

