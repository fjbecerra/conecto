package shopify

import (
	"encoding/json"
	"fmt"
)

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