package factories

import (
	"conecto/core"
	"conecto/core/checkpoint"
	"conecto/core/checkpoint/memory"
	"conecto/core/engines"
	"conecto/core/sinks"
	"conecto/core/sinks/codecs"
	"conecto/core/sinks/rdbs"
	"context"
	"math/rand/v2"
	"testing"
	"time"
)

func TestMockedFbAdInsightPipeline(t *testing.T) {
	dataStore := []map[string]interface{}{}
	stateStore:= make(map[string]core.State)
	stateStore["pipeline-test"] = State{
		Cursor : core.Cursor[]
	}
	pipelineTest := PipelineTest {
		ConfigPath: "./testdata/fb_ad_insights/ad_insight_pipeline.json",
		DataStore: dataStore,
		StateStore: map[string]checkpoint.StateStore{

		}

	}
	pipeline := BuildPipelineTest(pipelineTest)
	ctx :=context.Background()
	error:= pipeline.Run(ctx)
	if error != nil {
		t.Error(error.Error())
	}
	if len(dataStore) != 4 {
		t.Errorf("number of record expected is 4, returned: %d", len(dataStore))
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

type PipelineTest struct {
	ConfigPath string
	DataStore []map[string]interface{}
	StateStore map[string]core.State
}

func BuildPipelineTest(pipelineTest PipelineTest) engines.PipelineRunner{
	seed := time.Now().UnixNano()

	r := rand.New(rand.NewPCG(
		uint64(seed),
		uint64(seed>>1),
	))
	config := LoadConfigPipeline(pipelineTest.ConfigPath)
    connector := NewConnector(config.ConnectorConfig,r).Build()
	tranformer := NewTransform(config.TransformersConfig, config.AdditionalConfig).Build()
	codec:= codecs.JSONCodec{}
	adapter := rdbs.PostgresAdapter{Codec: &codec}
	sinkMemory := sinks.NewMemorySink(&pipelineTest.DataStore, &adapter)
	sink := engines.SinkEngine {
		Sink: sinkMemory,
		BatchSize: 10,
	}

	stateStore:= memory.MemoryStateStore{
		Store: pipelineTest.StateStore,
	}
	return &engines.Pipeline{
		ID : "pipeline-test",
		ConnectorEngine: &connector,
		SinkEngine:   &sink,
		Transformer: tranformer,
		Settings: engines.Settings{
			BufferSize: 10,
		},
		StateStore: &stateStore,
	}
}

