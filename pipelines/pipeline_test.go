package pipelines

import (
	"context"
	"testing"
)

func TestMockedFbAdInsightPipeline(t *testing.T) {

	registry := NewRegistryPipeline()

   	pipeline := registry.Factories["mockedFbAdInsight"]()

	ctx :=context.Background()
	pipeline.Run(ctx)
	data := pipeline.TestResult()
	if len(data) != 4 {
		t.Errorf("shoud return 4 records, recieved %d", len(data))
	}
}
