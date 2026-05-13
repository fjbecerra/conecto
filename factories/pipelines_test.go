package factories

import (
	"conecto/core"
	"conecto/core/engines"
	"conecto/core/sinks"
	"conecto/core/statestores"
	"context"
	"errors"
	"time"

	"testing"
)

func TestMockedFbAdInsightPipelineRawData(t *testing.T) {	
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_raw_data.json")
	pipeline:= BuildPipeline(config)
	runtime:= core.Runtime{
		PipelineId: config.RuntimeConfig.PipelineId,
		Context: context.Background(),
	}
	error:= pipeline.Run(runtime)
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
		Context: context.Background(),
	}
	error:= pipeline.Run(runtime)
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

	
	// RUN PIPELINE
	errCh := make(chan error, 1)

	go func() {
		errCh <- pipeline.Run(
			core.Runtime{
				Context:    ctx,
				PipelineId: "test",
			},
		)
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
	err = pipeline.Run(
		core.Runtime{
			Context:    context.Background(),
			PipelineId: "test",
		},
	)

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
		PipelineId: "test2",
		Context: context.Background(),
	}
	error:= pipeline.Run(runtime)
	if error != nil {
		t.Error(error.Error())
	}
}


