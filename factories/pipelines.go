package factories

import (
	"conecto/core/engines"
	"conecto/core/pipeline"
)

func BuildPipeline(config ConfigPipeline) pipeline.Pipeline{
	
	random:= &RandomImpl{}
	connections:= NewSource(config.SourcesConfig).Build()
	connectorRunnable := NewConnector(config.ConnectorConfig, random, connections).Build()
	transform := NewTransform(config.TransformersConfig, config.FieldsSpecsConfig, config.RuntimeConfig).Build()
	stateStore := NewStateStore(config.RuntimeConfig.StateStoreConfig, connections).Build()
	sinkCommiter := NewSink(config.SinkConfig, config.FieldsSpecsConfig, random, stateStore, connections).Build()
	
	engine := engines.Engine {
		ConnectorRunnable: connectorRunnable,
        Transformer: transform,
		SinkCommiter: sinkCommiter,
	}

	return pipeline.Pipeline{
		ID: config.ID,
		Engine: &engine,
        StateStore: stateStore,
	}
}