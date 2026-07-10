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
func TestFbAdInsightPipelineIntegrationTest(t *testing.T) {
	context := context.Background()
	config:= factories.LoadConfigPipeline("./testdata/ad_insight_pipeline_to_postgres.json")
	pipeline:= factories.BuildPipeline(config)
	credentialService:= pipeline.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*_http.HttpConnector).Provider.Client.CredentialService.(*credentials.AESGCMCredentialService)
	
	credential := credentials.Credential{
		Type: "oauth2",

		Data: map[string]string{
			"access_token": "shpat_xxxxx",
			"refresh_token": "xxxx",
		},

		Expiry: nil,
	}
	er:= credentialService.Save(context, config.ID, credential)
	if er != nil {
		t.Error(er.Error())
	}	
	error:= pipeline.Run(context)
	if error != nil {
		t.Error(error.Error())
	}
}


