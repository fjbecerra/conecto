package pipelines

import (
	"conecto/core"
	"conecto/core/sinks"
	"conecto/factories"
	"context"
	"testing"
)

func TestMockedFbAdInsightPipeline(t *testing.T) {
	store := []core.Record{}
	sinkMemory := sinks.NewMemorySink(&store)
	pipeline := BuildPipelineTest("./testdata/fb_ad_insights/ad_insight_pipeline.json",sinkMemory)
	ctx :=context.Background()
	pipeline.Run(ctx)
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


func BuildPipelineTest(configPath string ,sink *sinks.SinkMemory) PipelineRunner{
	config := core.LoadConfigPipeline(configPath)
    source := factories.NewSource(config.SourceConfig).Build()
	current := factories.NewTransform(source, config.TransformsConfig, config.AdditionalConfigs).Build()
	return &Pipeline[core.Record]{
		source: current,
		sink:   sink,
	}
}
