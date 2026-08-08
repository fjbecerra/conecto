package shopify

type ShopifyResourceConfig struct{
	AppURL       string   `json:"app_url"`
    ClientID     string   `json:"client_id"`
    ClientSecret string   `json:"client_secret"`
    Scopes       []string `json:"scopes"`
}

type ShopifyConnectorConfig struct{
	Name string `json:"name"`
	Query string `json:"query"`
}