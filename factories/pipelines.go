package factories

import (
	"conecto/auth/credentials"
	"conecto/core/engines"
	"conecto/core/pipelines"
	"conecto/core/retry"
	"conecto/core/statestores"
)

type Pipeline struct {
	connections  Connections
	random       retry.Random
	stateStore 	 statestores.StateStore
	credentialService credentials.CredentialService
	config 		 PipelineConfig
}

func NewPipeline(
	connections  Connections,
	random       retry.Random,
	stateStore 	 statestores.StateStore,
	credentialService credentials.CredentialService,
	config PipelineConfig,
	) *Pipeline{
	return &Pipeline{
		connections: connections,
		random: random,
		stateStore: stateStore,
		credentialService: credentialService,
		config: config,
	}
}

func (p *Pipeline) Build() pipelines.Pipeline{
	streams := []pipelines.Stream{}
	for _, streamConfig := range p.config.ConnectorConfig.StreamsConfig {
		connector := NewConnector(
			p.config.ConnectorConfig, 
			streamConfig, 
			p.random, 
			p.credentialService,
		).Build()
		transform := NewTransform(
			streamConfig.TransformersConfig, 
			streamConfig.FieldsSpecsConfig,
		).Build()
		sinkCommiter := NewSink(
			p.config.SinkConfig, 
			streamConfig.FieldsSpecsConfig, 
			p.random, 
			p.stateStore,
			p.connections, 
			streamConfig.DestinationConfig,
		).Build()
		

		engine := engines.Engine {
			ConnectorRunnable: connector,
			Transformer: transform,
			SinkCommiter: sinkCommiter,
		}
		stream := pipelines.Stream{
			Engine: &engine,
			StateStore: p.stateStore,
		}
		streams = append(streams, stream)
	}

	return pipelines.Pipeline{
		ID: p.config.ID,
		Streams: streams,
	}

}