package factories

import (
	"conecto/core/engines"
	"conecto/core/sinks"
	"conecto/core/sinks/committers"
	"conecto/core/sinks/statestores"
	"context"
	"errors"
	"sync"

	"testing"
)

func TestMockedFbAdInsightPipelineRawData(t *testing.T) {	
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_raw_data.json")
	pipeline:= BuildPipeline(config)
	runtime:= engines.Runtime{
		PipelineId: config.RuntimeConfig.PipelineId,
		Context: context.Background(),
	}
	error:= pipeline.Run(runtime)
	if error != nil {
		t.Error(error.Error())
	}

	memSink := pipeline.SinkEngine.Commiter.(*committers.MemoryCommitter).Sink.(*sinks.SinkMemory)
	
	if len(memSink.Mstore) != 4 {
		t.Errorf("number of record expected is 4, returned: %d", len(memSink.Mstore))
	}	
}

func TestMockedFbAdInsightPipelineFlattened(t *testing.T) {	
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_flattened_data.json")
	pipeline:= BuildPipeline(config)
	runtime:= engines.Runtime{
		PipelineId: config.RuntimeConfig.PipelineId,
		Context: context.Background(),
	}
	error:= pipeline.Run(runtime)
	if error != nil {
		t.Error(error.Error())
	}

	memSink := pipeline.SinkEngine.Commiter.(*committers.MemoryCommitter).Sink.(*sinks.SinkMemory)
	
	if len(memSink.Mstore) != 4 {
		t.Errorf("number of record expected is 4, returned: %d", len(memSink.Mstore))
	}	
}

func TestPipeline_CancelAndResume(t *testing.T) {
	
	ctx, cancel := context.WithCancel(context.Background())
	cfg := LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_test_pipeline_flattened_data.json")
	pipeline := BuildPipeline(cfg)

	// ACCESS IN-MEMORY COMPONENTS
	sink := pipeline.SinkEngine.Commiter.(*committers.MemoryCommitter).Sink.(*sinks.SinkMemory)
	store := pipeline.StateStore.(*statestores.MemoryStateStore)

	// CONTROL SIGNALS
	firstBatchDone := make(chan struct{})
	cancelOnce := sync.Once{}

	// BLOCK AFTER FIRST FLUSH
	sink.OnWrite = func() {

		cancelOnce.Do(func() {
			close(firstBatchDone)
		})
	}

	// RUN PIPELINE
	errCh := make(chan error, 1)

	go func() {
		errCh <- pipeline.Run(engines.Runtime{
			Context:    ctx,
			PipelineId: "test",
		})
	}()

	// WAIT UNTIL FIRST BATCH IS PROCESSED
	<-firstBatchDone

	// CANCEL PIPELINE
	cancel()

	err := <-errCh

	if err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}

	// VERIFY CHECKPOINT WAS SAVED
	state, err := store.Load(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	if state.Cursor == nil {
		t.Fatal("expected cursor to be persisted")
	}

	// RESTART PIPELINE
	ctx2 := context.Background()

	err = pipeline.Run(engines.Runtime{
		Context:    ctx2,
		PipelineId: "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	// VERIFY NO DUPLICATION ON RESUME
	if len(sink.Mstore) != 4 {
		t.Fatalf("expected 4 records, got %d", len(sink.Mstore))
	}
}



//todo run containers when integrations tests run
//this tests depends on postgres container. 
func TestFbAdInsightPipelineIntegrationTest(t *testing.T) {
	config:= LoadConfigPipeline("./testdata/fb_ad_insights/ad_insight_pipeline_with_db.json")
	pipeline:= BuildPipeline(config)
	runtime:= engines.Runtime{
		PipelineId: config.RuntimeConfig.PipelineId,
		Context: context.Background(),
	}
	error:= pipeline.Run(runtime)
	if error != nil {
		t.Error(error.Error())
	}
}


