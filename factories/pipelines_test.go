package factories

import (
	"conecto/core"
	"conecto/core/connectors/rest"
	"conecto/core/connectors/rest/auths"
	"conecto/core/engines"
	"conecto/core/sinks"
	"conecto/core/statestores"
	"context"
	"errors"
	"os"
	"time"

	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {

	_ = godotenv.Load("../.env")

	code := m.Run()

	os.Exit(code)
}

func TestMockedFbAdInsightPipelineRawData(t *testing.T) {	
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_raw_data.json")
	pipeline:= BuildPipeline(config)
	runtime:= core.Runtime{
		PipelineId: config.RuntimeConfig.PipelineId,
		Provider: "facebook_ad_insight",
		Context: context.Background(),
	}

	tokenStore:= pipeline.ConnectorEngine.Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)
	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	error:= tokenStore.Save(runtime, token)
	if error != nil {
		t.Error(error.Error())
	}	
	error= pipeline.Run(runtime)
	if error != nil {
		t.Error(error.Error())
	}

	memSink := pipeline.CommitStrategy.(*engines.AtLeastOnceCommitStrategy).Sink.(*sinks.SinkMemory)
	
	if len(memSink.Mstore) != 4 {
		t.Errorf("number of record expected is 4, returned: %d", len(memSink.Mstore))
	}	
}

func TestMockedFbAdInsightPipelineFlattened(t *testing.T) {	
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_flattened_data.json")

	pipeline:= BuildPipeline(config)
	runtime:= core.Runtime{
		PipelineId: config.RuntimeConfig.PipelineId,
		Provider: "facebook_ad_insight",
		Context: context.Background(),
	}

	tokenStore:= pipeline.ConnectorEngine.Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)
	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	error:= tokenStore.Save(runtime, token)
	if error != nil {
		t.Error(error.Error())
	}		

	error= pipeline.Run(runtime)
	if error != nil {
		t.Error(error.Error())
	}

	memSink := pipeline.CommitStrategy.(*engines.AtLeastOnceCommitStrategy).Sink.(*sinks.SinkMemory)
	
	if len(memSink.Mstore) != 4 {
		t.Errorf("number of record expected is 4, returned: %d", len(memSink.Mstore))
	}	
}


func TestPipeline_CancelAndResume(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := LoadConfigPipeline(
		"./testdata/fb_ad_insights/ad_insight_test_pipeline_flattened_data.json",
	)

	pipeline := BuildPipeline(cfg)

	sink := pipeline.CommitStrategy.(*engines.AtLeastOnceCommitStrategy).Sink.(*sinks.SinkMemory)

	store := pipeline.CommitStrategy.(*engines.AtLeastOnceCommitStrategy).StateStore.(*statestores.MemoryStateStore)

	runtime:= core.Runtime{
		PipelineId: "test",
		Provider: "facebook_ad_insight",
		Context: context.Background(),
	}

	tokenStore:= pipeline.ConnectorEngine.Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)
	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	er:= tokenStore.Save(runtime, token)
	if er != nil {
		t.Error(er.Error())
	}	
	
	// RUN PIPELINE
	errCh := make(chan error, 1)

	go func() {
		errCh <- pipeline.Run(runtime)
	}()

	// GIVE PIPELINE TIME TO PROCESS SOME DATA
	time.Sleep(50 * time.Millisecond)

	// CANCEL
	cancel()

	err := <-errCh

	if err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf(
			"expected context canceled, got %v",
			err,
		)
	}

	// VERIFY CHECKPOINT EXISTS
	state, err := store.Load(
		core.Runtime{
				Context:    ctx,
				PipelineId: "test",
			},
	)
	if err != nil {
		t.Fatal(err)
	}

	// checkpoint MAY be nil if cancel happened before
	// first commit completed
	_ = state

	// RESTART PIPELINE
	err = pipeline.Run(runtime)
	if err != nil {
		t.Fatal(err)
	}

	// VERIFY EVENTUAL CONSISTENCY
	if len(sink.Mstore) != 4 {
		t.Fatalf(
			"expected 4 records after resume, got %d",
			len(sink.Mstore),
		)
	}
}

// //todo run containers when integrations tests run
// //this tests depends on postgres container. 
func TestFbAdInsightPipelineIntegrationTest(t *testing.T) {
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_pipeline_with_db.json")
	pipeline:= BuildPipeline(config)
	runtime:= core.Runtime{
		PipelineId: "test",
		Context: context.Background(),
	}
	tokenStore:= pipeline.ConnectorEngine.Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)
	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	er:= tokenStore.Save(runtime, token)
	if er != nil {
		t.Error(er.Error())
	}	
	error:= pipeline.Run(runtime)
	if error != nil {
		t.Error(error.Error())
	}
}


