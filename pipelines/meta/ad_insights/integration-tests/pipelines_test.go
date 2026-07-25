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
func TestFbAdInsightPipelineIntegrationTest(t *testing.T) {
	context := context.Background()
	config := factories.LoadConfigPipeline("./testdata/ad_insight_pipeline_to_postgres.json")
	pipeline := factories.BuildPipeline(config)
	connection := connections.Connection{
		ID: "test-connection-id-123-ad-insights",
	}
	for _, stream := range pipeline.Streams {
		credentialService := stream.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*api.HttpConnector).Provider.Client.CredentialService.(*credentials.AESGCMCredentialService)

		credential := credentials.Credential{
			Type: "oauth2",

			Data: map[string]string{
				"access_token":  "shpat_xxxxx",
				"refresh_token": "xxxx",
			},

			Expiry: nil,
		}
		er := credentialService.Save(context, connection, credential)
		if er != nil {
			t.Error(er.Error())
		}
		error := stream.Run(context, connection)
		if error != nil {
			t.Error(error.Error())
		}
	}

}
