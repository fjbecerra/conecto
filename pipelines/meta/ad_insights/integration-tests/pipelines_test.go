package factories

import (
	"conecto/connectors/rest"
	"conecto/connectors/rest/auths"
	"conecto/core/engines"
	"conecto/factories"
	"context"
	"os"
	"testing"
	"time"
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
	tokenStore:= pipeline.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)
	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	er:= tokenStore.Save(context,config.ID, token)
	if er != nil {
		t.Error(er.Error())
	}	
	error:= pipeline.Run(context)
	if error != nil {
		t.Error(error.Error())
	}
}


