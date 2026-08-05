package shopify

import (
	"conecto/auth/connections"
	"fmt"
)

type ShopifyEndpointProvider struct {}

func (p *ShopifyEndpointProvider) Apply(connection connections.Connection) string {
	shop := connection.Metadata["shop"]
	apiVersion := connection.Metadata["api_version"]
	return fmt.Sprintf("https://%s.myshopify.com/admin/api/%s/graphql.json", shop, apiVersion)
}