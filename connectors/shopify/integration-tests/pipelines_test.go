package factories

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"conecto/connectors/api"
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
	config := factories.LoadConfigPipeline("./testdata/orders_pipeline_to_postgres.json")
	pipeline := factories.BuildPipeline(config)
	connection := connections.Connection{
		ID: "test-connection-id-123-shopify",
	}
	for _, stream := range pipeline.Streams {
		credentialService := stream.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*api.HttpConnector).Provider.Client.CredentialService.(*credentials.AESGCMCredentialService)

		credential := credentials.Credential{
			Type: "oauth2",

			Data: map[string]string{
				"X-Shopify-Access-Token": "shpat_xxxxx",
				"refresh_token":          "xxx",
			},

			Expiry: nil,
		}
		error := credentialService.Save(context, connection, credential)
		if error != nil {
			t.Error(error.Error())
		}
		error = stream.Run(context, connection)
		if error != nil {
			t.Error(error.Error())
		}
	}

}
