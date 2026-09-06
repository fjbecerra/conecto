package shopify

type ShopifyResourceConfig struct{
	AppURL       string   `json:"app_url"`
    ClientID     string   `json:"client_id"`
    ClientSecret string   `json:"client_secret"`
    Scopes       []string `json:"scopes"`
}

type ShopifyConnectorConfig struct{
	BatchSize int `json:"batch_size"`
	BackfillLastNDays int `json:"backfill_last_n_days"`
	IncrementalLastNDays int `json:"incremental_last_n_days"`
}