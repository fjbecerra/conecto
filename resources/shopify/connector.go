package shopify

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/retry"
	"conecto/core/statestores"
	"conecto/resources/base/api"
	"conecto/resources/base/api/graphql"
	"conecto/shared/config"
	"conecto/stores/credentials"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)


type ShopifyConnector struct {
	name			  	string
	fieldSpecs			config.FieldsSpecs
    httpClient        	*api.HttpClient
    credentialService 	credentials.CredentialService
	retryExecutor 		*retry.Executor
	cfg 				ShopifyConnectorConfig
}

func CreateShopifyConnector(shopifyConnector ShopifyConnector) engines.ConnectorRunnable {
	provider:= &api.HeaderProvider{
		Name: "X-Shopify-Access-Token",
	}
	
	builder:= &graphql.GraphQLRequestBuilder{
		EndpointProvider:  &ShopifyEndpointProvider{},
		Query: buildQuery(shopifyConnector.name, shopifyConnector.fieldSpecs, shopifyConnector.cfg.BatchSize),
		VariableCursorKey: "after",
		SyncRequestProvider: &ShopifySyncRequestProvider{
			backfillLastNDays: shopifyConnector.cfg.Sync.BackfillLastNDays,
			watermarkLastNDays: shopifyConnector.cfg.Sync.IncrementalLastNDays,
		},

	}

	dataExtractor := &graphql.GraphQLDataExtractor{
		Path: fmt.Sprintf("data.%s.edges", shopifyConnector.name),
	}

	cursorExtractor := &graphql.GraphQLCursorExtractor{
		HasMorePath: fmt.Sprintf("data.%s.pageInfo.hasNextPage", shopifyConnector.name),
		CursorPath:  fmt.Sprintf("data.%s.pageInfo.endCursor", shopifyConnector.name),
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

func  buildQuery(name string, fieldSpecs config.FieldsSpecs, batchSize int) string{
	fields := []string{}
	for _, fieldSpec := range fieldSpecs {
		fields = append(fields, fieldSpec.Path)
	}
	query := `query %sPage($after: String, $query: String) {
				%s(first: %d, after: $after, query: $query) {
					edges {
						node {
							%s
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}`
	return fmt.Sprintf(
		query,
		name,       
		name,      
		batchSize,         
		strings.Join(fields, " "),
	)
}

type ShopifySyncRequestProvider struct {
	backfillLastNDays int
	watermarkLastNDays int
}
func (s *ShopifySyncRequestProvider) Apply(syncState statestores.SyncState) map[string]any {
	if(syncState == nil){
		if(s.backfillLastNDays > 0){
			t := time.Now()
			from := t.AddDate(0, 0, -s.backfillLastNDays)
			return map[string]any{
				"query": fmt.Sprintf("updated_at>%s", from.UTC().Format(time.RFC3339)),
			}
		}
		return nil	
	}	
	if(syncState["watermark"] != ""){
		t, _ := time.Parse(time.RFC3339, syncState["watermark"])
		from := t.AddDate(0, 0, s.watermarkLastNDays)
		return map[string]any{
				"query": fmt.Sprintf("updated_at>%s", from.UTC().Format(time.RFC3339)),
		}
	}
	return nil
	
}