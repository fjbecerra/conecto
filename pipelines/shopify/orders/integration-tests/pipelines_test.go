package factories

import (
	"conecto/auth/credentials"
	"conecto/connectors/_http"
	"conecto/core/engines"
	"conecto/factories"
	"context"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Setenv("TOKEN_ENCRYPTION_KEY_V1", "z1Wbj51mq1GmIpmwAfLv9X5oSekOYEsC/9YXhOCuKjU=")
	code := m.Run()
	os.Exit(code)
}

// //todo run containers when integrations tests run
// //this tests depends on postgres container. 
func TestShopifyPipelineIntegrationTest(t *testing.T) {
	context := context.Background()
	config:= factories.LoadConfigPipeline("./testdata/orders_pipeline_to_postgres.json")
	pipeline:= factories.BuildPipeline(config)
	credentialService:= pipeline.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*_http.HttpConnector).Provider.Client.CredentialService.(*credentials.AESGCMCredentialService)
	
	credential := credentials.Credential{
		Type: "oauth2",

		Data: map[string]string{
			"X-Shopify-Access-Token": "shpat_xxxxx",
			"refresh_token": "xxxx",
		},

		Expiry: nil,
	}
	error:= credentialService.Save(context, config.ID, credential)
	if error != nil {
		t.Error(error.Error())
	}
	error= pipeline.Run(context)
	if error != nil {
		t.Error(error.Error())
	}
}


