package factories

import (
	"conecto/core/engines"
	"conecto/core/streams"
)

func BuildStreams(conecto Conecto, config PipelineConfig) []streams.Stream{
	
	// random:= &RandomImpl{}
	// sources:= NewSource(config.SourcesConfig).Build()
	
	streams := []streams.Stream{}
	for _, streamConfig := range config.StreamsConfig {
		connector := NewConnector(config.ConnectorConfig, streamConfig, conecto.random, conecto.connections).Build()
		transform := NewTransform(streamConfig.TransformersConfig, streamConfig.FieldsSpecsConfig).Build()
		//stateStore := NewStateStore(config.RuntimeConfig.StateStoreConfig, sources).Build()
		sinkCommiter := NewSink(config.SinkConfig, streamConfig.FieldsSpecsConfig, conecto.random, conecto.stateStore, conecto.connections).Build()
		

		engine := engines.Engine {
			ConnectorRunnable: connector,
			Transformer: transform,
			SinkCommiter: sinkCommiter,
		}
		stream := streams.Stream{
			Engine: &engine,
			StateStore: conecto.stateStore,
		}
		streams = append(streams, stream)
	}

	return streams
}