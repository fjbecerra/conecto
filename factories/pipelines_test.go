package factories

import (
	"conecto/connectors/rest"
	"conecto/connectors/rest/auths"
	"conecto/core/engines"
	"conecto/sinks/memory"
	"conecto/states"
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
	context := context.Background()
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_raw_data.json")
	pipeline:= BuildPipeline(config)
	
	tokenStore:= pipeline.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)
	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	error:= tokenStore.Save(context, config.ID, token)
	if error != nil {
		t.Error(error.Error())
	}	
	error= pipeline.Run(context)
	if error != nil {
		t.Error(error.Error())
	}

	memSink := pipeline.Engine.SinkCommiter.(*engines.SinkEngine).Sink.(*memory.SinkMemory)
	
	if len(memSink.Mstore) != 4 {
		t.Errorf("number of record expected is 4, returned: %d", len(memSink.Mstore))
	}	
}

func TestMockedFbAdInsightPipelineFlattened(t *testing.T) {	
	context := context.Background()
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_flattened_data.json")

	pipeline:= BuildPipeline(config)

	tokenStore:= pipeline.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)

	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	error:= tokenStore.Save(context, config.ID, token)
	if error != nil {
		t.Error(error.Error())
	}		

	error= pipeline.Run(context)
	if error != nil {
		t.Error(error.Error())
	}

	memSink := pipeline.Engine.SinkCommiter.(*engines.SinkEngine).Sink.(*memory.SinkMemory)
	
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

	sink := pipeline.Engine.SinkCommiter.(*engines.SinkEngine).Sink.(*memory.SinkMemory)

	store := pipeline.Engine.SinkCommiter.(*engines.SinkEngine).StateStore.(*states.MemoryStateStore)

	tokenStore:= pipeline.Engine.ConnectorRunnable.(*engines.ConnectorEngine).Connector.(*rest.RESTConnector).Provider.RestClient.TokenStore.(*auths.AESGCMTokenStore)

	token:= auths.Token{
		AccessToken : "any-token",
		RefreshToken: "any-refresh-token",
		Expiry: time.Now(),
	}
	er:= tokenStore.Save(ctx,cfg.ID, token)
	if er != nil {
		t.Error(er.Error())
	}	
	
	// RUN PIPELINE
	errCh := make(chan error, 1)

	go func() {
		errCh <- pipeline.Run(ctx)
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
	state, err := store.Load(ctx, cfg.ID)
	if err != nil {
		t.Fatal(err)
	}

	// checkpoint MAY be nil if cancel happened before
	// first commit completed
	_ = state

	// RESTART PIPELINE
	err = pipeline.Run(ctx)
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
	context := context.Background()
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_pipeline_with_db.json")
	pipeline:= BuildPipeline(config)
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


