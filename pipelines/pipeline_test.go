package pipelines

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/sinks"
	"conecto/core/sinks/codecs"
	"conecto/core/sinks/rdbs"
	"conecto/factories"
	"context"
	"testing"
)

func TestMockedFbAdInsightPipeline(t *testing.T) {
	store := []map[string]interface{}{}
	codec:= codecs.JSONCodec{}
	adapter := rdbs.PostgresAdapter{Codec: &codec}
	sinkMemory := sinks.NewMemorySink(&store, &adapter)
	pipeline := BuildPipelineTest("./testdata/fb_ad_insights/ad_insight_pipeline.json",sinkMemory)
	ctx :=context.Background()
	error:= pipeline.Run(ctx)
	if error != nil {
		t.Error(error.Error())
	}
	if len(store) != 4 {
		t.Errorf("number of record expected is 4, returned: %d", len(store))
	}
}

//todo run containers when integrations tests run
//this tests depends on postgres container. 
func TestFbAdInsightPipelineIntegrationTest(t *testing.T) {

	registry := NewRegistryPipeline()

   	pipeline := registry.Factories["fbAdInsight"]()

	ctx :=context.Background()
	error:=pipeline.Run(ctx)
	
	if error != nil {
		t.Error(error.Error())
	}
}


func BuildPipelineTest(configPath string ,sinkMemory *sinks.SinkMemory) PipelineRunner{
	config := core.LoadConfigPipeline(configPath)
    connector := factories.NewConnector(config.ConnectorConfig).Build()
	tranformer := factories.NewTransform(config.TransformersConfig, config.AdditionalConfig).Build()
	sink := engines.SinkEngine {
		Sink: sinkMemory,
		BatchSize: 10,
	}
	return &Pipeline{
		connectorEngine: &connector,
		sinkEngine:   &sink,
		transformer: tranformer,
		bufferSize: 10,
	}
}
