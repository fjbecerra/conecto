package factories

import (
	"conecto/auth/connections"
	"conecto/core/engines"
	"conecto/core/pipeline"
	
)

func BuildPipeline(config ConfigPipeline) pipeline.Pipeline{
	
	random:= &RandomImpl{}
	sources:= NewSource(config.SourcesConfig).Build()
	connectorRunnable := NewConnector(config.ConnectorConfig, random, sources).Build()
	transform := NewTransform(config.TransformersConfig, config.FieldsSpecsConfig, config.RuntimeConfig).Build()
	stateStore := NewStateStore(config.RuntimeConfig.StateStoreConfig, sources).Build()
	sinkCommiter := NewSink(config.SinkConfig, config.FieldsSpecsConfig, random, stateStore, sources).Build()
	

	engine := engines.Engine {
		ConnectorRunnable: connectorRunnable,
        Transformer: transform,
		SinkCommiter: sinkCommiter,
	}

	// connectionStore := connections.NewMemoryStore()
	// connectionStore.Get(ctx, config.ID)

	return pipeline.Pipeline{
		Connection: connections.Connection{
			ID: config.ID,
			TenantID: "agency_x",
			Provider: "shopify",
			ExternalID: "shop1",
			Metadata: map[string]any{
				"shop": "shop1",
			},
			Status: "active",
		},
		Engine: &engine,
        StateStore: stateStore,
	}
}